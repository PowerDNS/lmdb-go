package lmdb

import "errors"

// v10Available is false on Windows: upstream LMDB 1.0 is broken there, so the
// 1.0 tree is excluded from the build (//go:build !windows) and only the 0.9
// engine exists. Requesting V10, or opening a database in the 1.0 format,
// fails with errV10Unavailable.
const v10Available = false

var errV10Unavailable = errors.New("lmdb: the LMDB 1.0 engine is not supported on Windows")
