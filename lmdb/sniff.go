package lmdb

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Meta page layout facts (see mdb.c of both streams): both formats start the
// MDB_meta (mm_magic, mm_version, ...) right after the page header, whose
// size is 16 bytes in 0.9 and 24 in 1.0 (1.0 added an 8-byte mh_txnid).
// mm_magic is 0xBEEFC0DE in both; mm_version (MDB_DATA_VERSION) is 1 in 0.9
// and 3 in 1.0. Version 2 was written by lmdb-js prerelease snapshots and is
// compatible with neither. The fields are written in the platform's native
// byte order.
const sniffMagic = 0xBEEFC0DE

// ErrFormatUnsupported is returned when a database file carries a valid LMDB
// magic but data format version 2, which belongs to neither bundled engine.
// It is written by lmdb-js, whose vendored LMDB prerelease snapshot defines
// MDB_DATA_VERSION 2 (verified against kriszyp/lmdb-js
// dependencies/lmdb/libraries/liblmdb/mdb.c).
var ErrFormatUnsupported = errors.New("lmdb: unsupported database format (written by an lmdb-js prerelease)")

// SniffFormat detects which LMDB version matches the on-disk format of the
// database at path, without opening the environment. Opening with the wrong
// engine to find out is not safe: both engines set up (and may rewrite) the
// lock file before reading the data file's header — mdb_env_open calls
// mdb_env_setup_locks before mdb_env_open2/mdb_env_read_header — so a failed
// probe can race a concurrent correct-version open.
//
// Of flags, only NoSubdir is honored (path is the data file itself rather
// than a directory containing data.mdb).
//
// A missing or empty data file yields (FormatUnknown, nil): the environment
// does not exist yet and Open would create it. A file that is not an LMDB
// database yields an OpError wrapping Invalid; a valid LMDB file of an
// unknown data version yields an OpError wrapping VersionMismatch; the
// lmdb-js prerelease format yields ErrFormatUnsupported.
func SniffFormat(path string, flags uint) (LMDBVersion, error) {
	datafile := path
	if flags&NoSubdir == 0 {
		datafile = filepath.Join(path, "data.mdb")
	}

	f, err := os.Open(datafile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return FormatUnknown, nil
		}
		return FormatUnknown, err
	}
	defer f.Close()

	var buf [32]byte
	n, err := io.ReadFull(f, buf[:])
	if n == 0 && errors.Is(err, io.EOF) {
		return FormatUnknown, nil
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return FormatUnknown, &OpError{Op: "sniff", Errno: Invalid}
	}
	if err != nil {
		return FormatUnknown, err
	}

	var mmver uint32
	switch {
	case binary.NativeEndian.Uint32(buf[16:]) == sniffMagic: // 0.9 page header size
		mmver = binary.NativeEndian.Uint32(buf[20:])
	case binary.NativeEndian.Uint32(buf[24:]) == sniffMagic: // 1.0 page header size
		mmver = binary.NativeEndian.Uint32(buf[28:])
	default:
		return FormatUnknown, &OpError{Op: "sniff", Errno: Invalid}
	}

	switch mmver {
	case 1:
		return V09, nil
	case 3:
		return V10, nil
	case 2:
		return FormatUnknown, ErrFormatUnsupported
	default:
		// Valid magic, unknown data version (e.g. MDB_DEVEL builds).
		return FormatUnknown, &OpError{Op: "sniff", Errno: VersionMismatch}
	}
}
