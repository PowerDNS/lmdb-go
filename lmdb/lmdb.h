/* lmdb-go v2 lmdb.h shim.
 *
 * Two vendored LMDB streams live side by side in this package, each compiled
 * behind a generated rename header (rename_lmdb09.h / rename_lmdb10.h) that
 * prefixes every extern symbol (mdb09_* / mdb10_*) and defines LMDBGO_STREAM.
 *
 * Roles of this shim:
 *  - For the vendored engine sources (which define LMDBGO_STREAM via their
 *    rename header), selects the matching pristine upstream header.
 *  - For consumers (cgo preambles, lmdbgo.c; no LMDBGO_STREAM), presents the
 *    LMDB 1.0 header as the canonical API surface: shared MDB_* defines are
 *    identical between the streams (enforced by scripts/check-defines.sh) and
 *    MDB_val/MDB_stat/MDB_envinfo are layout-identical on 64-bit non-VL32
 *    builds (enforced by _Static_asserts in the per-stream check TUs).
 */
#ifndef LMDBGO_LMDB_H_SHIM
#define LMDBGO_LMDB_H_SHIM

#ifdef MDB_VL32
# error "lmdb-go does not support MDB_VL32 builds"
#endif

#if !defined(LMDBGO_STREAM)
# include "lmdb_lmdb10.h"
#elif LMDBGO_STREAM == 9
# include "lmdb_lmdb09.h"
#elif LMDBGO_STREAM == 10
# include "lmdb_lmdb10.h"
#else
# error "invalid LMDBGO_STREAM (must be 9 or 10)"
#endif

#endif /* LMDBGO_LMDB_H_SHIM */
