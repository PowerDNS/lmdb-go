package lmdb

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// createEnvFile creates a database of the given version in dir and returns
// after closing it.
func createEnvFile(t *testing.T, dir string, ver Format, flags uint) {
	t.Helper()
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if ver.valid() {
		if err := env.SetFormat(ver); err != nil {
			t.Fatal(err)
		}
	}
	if err := env.Open(dir, flags, 0644); err != nil {
		t.Fatal(err)
	}
	err = env.Update(func(txn *Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		return txn.Put(db, []byte("sniffkey"), []byte("sniffval"), 0)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSniffFormat_v09(t *testing.T) {
	dir := t.TempDir()
	createEnvFile(t, dir, V09, 0)
	ver, err := SniffFormat(dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ver != V09 {
		t.Errorf("sniffed %v (!= %v)", ver, V09)
	}
}

func TestSniffFormat_v10(t *testing.T) {
	dir := t.TempDir()
	createEnvFile(t, dir, V10, 0)
	ver, err := SniffFormat(dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ver != V10 {
		t.Errorf("sniffed %v (!= %v)", ver, V10)
	}
}

func TestSniffFormat_noSubdir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.file")
	createEnvFile(t, path, V09, NoSubdir)
	ver, err := SniffFormat(path, NoSubdir)
	if err != nil {
		t.Fatal(err)
	}
	if ver != V09 {
		t.Errorf("sniffed %v (!= %v)", ver, V09)
	}
}

func TestSniffFormat_missing(t *testing.T) {
	// Missing directory, missing data file, and empty data file are all
	// "no database yet": (FormatUnknown, nil).
	for _, tc := range []struct {
		name string
		prep func(t *testing.T) (path string, flags uint)
	}{
		{"missingDir", func(t *testing.T) (string, uint) {
			return filepath.Join(t.TempDir(), "does-not-exist"), 0
		}},
		{"emptyDir", func(t *testing.T) (string, uint) {
			return t.TempDir(), 0
		}},
		{"emptyFile", func(t *testing.T) (string, uint) {
			dir := t.TempDir()
			f := filepath.Join(dir, "data.mdb")
			if err := os.WriteFile(f, nil, 0644); err != nil {
				t.Fatal(err)
			}
			return dir, 0
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path, flags := tc.prep(t)
			ver, err := SniffFormat(path, flags)
			if err != nil {
				t.Fatal(err)
			}
			if ver != FormatUnknown {
				t.Errorf("sniffed %v (!= FormatUnknown)", ver)
			}
		})
	}
}

// sniffFixture writes a 32-byte fake data file with the LMDB magic at the
// given offset and the data version right after it.
func sniffFixture(t *testing.T, magicOffset int, dataVersion uint32) string {
	t.Helper()
	var buf [32]byte
	binary.NativeEndian.PutUint32(buf[magicOffset:], sniffMagic)
	binary.NativeEndian.PutUint32(buf[magicOffset+4:], dataVersion)
	path := filepath.Join(t.TempDir(), "data.mdb")
	if err := os.WriteFile(path, buf[:], 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSniffFormat_fixtures(t *testing.T) {
	for _, tc := range []struct {
		name        string
		magicOffset int
		dataVersion uint32
		wantVer     Format
		wantErrno   Errno // 0 means no Errno expected
		wantJSErr   bool
	}{
		{"v09", 16, 1, V09, 0, false},
		{"v10", 24, 3, V10, 0, false},
		// lmdb-js prerelease snapshots (data version 2), either header size.
		{"lmdbjs09hdr", 16, 2, FormatUnknown, 0, true},
		{"lmdbjs10hdr", 24, 2, FormatUnknown, 0, true},
		// Valid magic, unknown data version (e.g. MDB_DEVEL builds).
		{"unknownVersion", 16, 7, FormatUnknown, VersionMismatch, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := sniffFixture(t, tc.magicOffset, tc.dataVersion)
			ver, err := SniffFormat(path, NoSubdir)
			if ver != tc.wantVer {
				t.Errorf("sniffed %v (!= %v)", ver, tc.wantVer)
			}
			switch {
			case tc.wantJSErr:
				if !errors.Is(err, ErrFormatUnsupported) {
					t.Errorf("expected ErrFormatUnsupported, got %v", err)
				}
			case tc.wantErrno != 0:
				if !IsErrno(err, tc.wantErrno) {
					t.Errorf("expected errno %v, got %v", tc.wantErrno, err)
				}
			case err != nil:
				t.Error(err)
			}
		})
	}
}

func TestSniffFormat_notLMDB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.mdb")
	if err := os.WriteFile(path, []byte("this is not an lmdb file at all, definitely"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := SniffFormat(path, NoSubdir)
	if !IsErrno(err, Invalid) {
		t.Errorf("expected Invalid, got %v", err)
	}

	// A short file that is not empty is also not a valid database.
	if err := os.WriteFile(path, []byte("short"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = SniffFormat(path, NoSubdir)
	if !IsErrno(err, Invalid) {
		t.Errorf("expected Invalid, got %v", err)
	}
}

// TestSniffFormat_openMismatch verifies that Open on a fake unsupported file
// fails with the sniffer's precise error rather than the engine's.
func TestSniffFormat_openMismatch(t *testing.T) {
	path := sniffFixture(t, 16, 2)
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	err = env.Open(path, NoSubdir, 0644)
	if !errors.Is(err, ErrFormatUnsupported) {
		t.Errorf("expected ErrFormatUnsupported, got %v", err)
	}
}
