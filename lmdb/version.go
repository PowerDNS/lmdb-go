package lmdb

/*
#include "lmdb.h"
#include "lmdbgo.h"
*/
import "C"

import (
	"errors"
	"os"
	"sync/atomic"
)

// Format identifies one of the two LMDB engines bundled in this package.
//
// Both engines are always linked into the binary. Each Env is driven by
// exactly one engine, selected when Open is called: an existing database is
// always opened with the engine matching its on-disk format (detected by
// SniffFormat), while a new database uses the explicitly requested version
// (Env.SetFormat) or the process default (SetDefaultFormat,
// DefaultFormatEnvVar, or V09).
type Format int

const (
	// FormatUnknown means no LMDB version has been determined (yet).
	FormatUnknown Format = 0
	// V09 selects the bundled LMDB 0.9.x engine (on-disk data format 1).
	V09 Format = 9
	// V10 selects the bundled LMDB 1.0.x engine (on-disk data format 3).
	// The two formats are mutually incompatible on disk.
	V10 Format = 10
)

func (v Format) valid() bool { return v == V09 || v == V10 }

// String returns "0.9", "1.0" or "unknown".
func (v Format) String() string {
	switch v {
	case V09:
		return "0.9"
	case V10:
		return "1.0"
	default:
		return "unknown"
	}
}

// DefaultFormatEnvVar is the environment variable consulted for the
// process-wide default LMDB version when SetDefaultFormat has not been
// called. Accepted values: "09", "9", "0.9", "10", "1.0".
const DefaultFormatEnvVar = "LMDBGO_DEFAULT_FORMAT"

// ErrInvalidFormat is returned when an Format value is not V09 or V10.
var ErrInvalidFormat = errors.New("lmdb: invalid LMDB version")

// ErrVersionUndetermined is returned by Env.EngineVersion before the engine
// driving the environment is known (no explicit request and not yet opened).
var ErrVersionUndetermined = errors.New("lmdb: LMDB version not determined yet")

// ErrFormatConflict is returned by Env.Open when the on-disk format of an
// existing database differs from the version explicitly requested with
// Env.SetFormat. lmdb-go never converts between formats; migrating
// requires a dump and reload (which can be done in a single process by
// opening a source and a destination Env of different versions).
var ErrFormatConflict = errors.New("lmdb: existing database format does not match the explicitly requested LMDB version")

var defaultVersion atomic.Int32

// SetDefaultFormat sets the process-wide default LMDB version used when
// Env.Open creates a new database and no version was requested on the Env.
// Passing FormatUnknown reverts to the built-in behavior (the
// DefaultFormatEnvVar environment variable, or V09).
//
// The default is V09: while LMDB 0.9 remains what distributions and other
// tools ship, databases created in the 0.9 format stay readable by them.
func SetDefaultFormat(v Format) error {
	switch v {
	case FormatUnknown, V09:
	case V10:
		if !v10Available {
			return errV10Unavailable
		}
	default:
		return ErrInvalidFormat
	}
	defaultVersion.Store(int32(v))
	return nil
}

// defaultFormat resolves the default version for new databases:
// SetDefaultFormat wins over DefaultFormatEnvVar, which wins over V09.
func defaultFormat() Format {
	if v := Format(defaultVersion.Load()); v.valid() {
		return v
	}
	switch os.Getenv(DefaultFormatEnvVar) {
	case "10", "1.0":
		if v10Available {
			return V10
		}
	case "09", "9", "0.9":
		return V09
	}
	return V09
}

// EnvVersion describes the exact LMDB engine version driving an Env.
type EnvVersion struct {
	Major   int
	Minor   int
	Patch   int
	Release string // human-readable upstream release string
}

// String returns the upstream release string.
func (v EnvVersion) String() string { return v.Release }

// SetFormat requests the LMDB version used if Open creates a NEW
// database. It must be called before Open.
//
// The request only applies to database creation: an existing database is
// always opened with the engine matching its on-disk format, and if that
// format differs from an explicitly requested version, Open fails with
// ErrFormatConflict rather than converting anything.
//
// Passing FormatUnknown clears the request.
func (env *Env) SetFormat(v Format) error {
	switch v {
	case FormatUnknown, V09:
	case V10:
		if !v10Available {
			return errV10Unavailable
		}
	default:
		return ErrInvalidFormat
	}

	env.closeLock.Lock()
	defer env.closeLock.Unlock()
	if env.opened || env.closed {
		return errors.New("lmdb: SetFormat must be called before Open")
	}
	env.reqVer = v
	return nil
}

// Format returns the LMDB version driving this environment. Before Open
// it returns the version requested with SetFormat, or FormatUnknown.
func (env *Env) Format() Format {
	env.closeLock.RLock()
	defer env.closeLock.RUnlock()
	if env.opened {
		return Format(env.ver)
	}
	return env.reqVer
}

// EngineVersion returns the exact version of the LMDB engine driving this
// environment. Before Open it reports the explicitly requested engine, if
// any; otherwise it returns ErrVersionUndetermined.
//
// The package-level Version and VersionString report the canonical (newest
// bundled) LMDB identity instead; EngineVersion is the per-environment truth.
func (env *Env) EngineVersion() (EnvVersion, error) {
	v := env.Format()
	if !v.valid() {
		return EnvVersion{}, ErrVersionUndetermined
	}
	var maj, min, pat C.int
	verstr := C.lmdbgo2_mdb_version(C.int(v), &maj, &min, &pat)
	return EnvVersion{
		Major:   int(maj),
		Minor:   int(min),
		Patch:   int(pat),
		Release: C.GoString(verstr),
	}, nil
}
