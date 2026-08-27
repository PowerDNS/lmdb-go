//go:build !windows

package lmdb

import "errors"

// v10Available reports whether the bundled LMDB 1.0 engine exists in this
// build. It is false only on Windows, where upstream LMDB 1.0 is broken and
// the 1.0 tree is excluded from the build.
const v10Available = true

// errV10Unavailable is never returned on this platform but must compile.
var errV10Unavailable = errors.New("lmdb: LMDB 1.0 engine is not available in this build")
