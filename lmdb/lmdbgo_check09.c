/* lmdbgo_check09.c
 *
 * Compile-time contract checks for the LMDB 0.9 stream. Defines no symbols.
 *
 * The rename header applies the mdb_* -> mdb09_* renames and sets
 * LMDBGO_STREAM, so the shim lmdb.h below resolves to the stream's own
 * pristine header with renamed declarations. Re-declaring the same symbols
 * via lmdbgo_dispatch.h (with the canonical prototypes the dispatch shims
 * use) makes any prototype drift between this engine and the canonical
 * surface a compile error in this TU.
 */
#include "rename_lmdb09.h"
#include "lmdb.h"
#include <stddef.h>

#define LMDBGO_SYM(name) mdb09_##name
#include "lmdbgo_dispatch.h"
#undef LMDBGO_SYM

/* lmdb-go v2 supports only 64-bit, non-MDB_VL32 targets: that is what makes
 * MDB_val/MDB_stat/MDB_envinfo layout-identical between the streams, so the
 * canonical (LMDB 1.0) header can serve both engines. */
_Static_assert(sizeof(size_t) == 8, "lmdb-go v2 requires a 64-bit target");
_Static_assert(sizeof(void *) == 8, "lmdb-go v2 requires 64-bit pointers");

_Static_assert(sizeof(MDB_val) == 16, "MDB_val layout contract");
_Static_assert(offsetof(MDB_val, mv_size) == 0, "MDB_val layout contract");
_Static_assert(offsetof(MDB_val, mv_data) == 8, "MDB_val layout contract");

_Static_assert(sizeof(MDB_stat) == 40, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_psize) == 0, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_depth) == 4, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_branch_pages) == 8, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_leaf_pages) == 16, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_overflow_pages) == 24, "MDB_stat layout contract");
_Static_assert(offsetof(MDB_stat, ms_entries) == 32, "MDB_stat layout contract");

_Static_assert(sizeof(MDB_envinfo) == 40, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_mapaddr) == 0, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_mapsize) == 8, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_last_pgno) == 16, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_last_txnid) == 24, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_maxreaders) == 32, "MDB_envinfo layout contract");
_Static_assert(offsetof(MDB_envinfo, me_numreaders) == 36, "MDB_envinfo layout contract");

/* Spot checks that shared constants carry the canonical values (full-surface
 * drift is checked by scripts/check-defines.sh). */
_Static_assert(MDB_PREV_MULTIPLE == 18, "MDB_cursor_op contract");
_Static_assert(MDB_APPENDDUP == 0x40000, "flag value contract");
_Static_assert(MDB_KEYEXIST == -30799, "errno range contract");
