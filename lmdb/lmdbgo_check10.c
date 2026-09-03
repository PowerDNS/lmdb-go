/* lmdbgo_check10.c
 *
 * Compile-time contract checks for the LMDB 1.0 stream. Defines no symbols.
 * See lmdbgo_check09.c for how the redeclaration trick works.
 */
#include "rename_lmdb10.h"
#include "lmdb.h"
#include <stddef.h>

#define LMDBGO_SYM(name) mdb10_##name
#include "lmdbgo_dispatch.h"
#undef LMDBGO_SYM

_Static_assert(sizeof(size_t) == 8, "lmdb-go v2 requires a 64-bit target");
_Static_assert(sizeof(void *) == 8, "lmdb-go v2 requires 64-bit pointers");

/* On non-MDB_VL32 builds mdb_size_t must be size_t itself; the dispatch
 * declarations above already prove this per parameter, this documents it. */
_Static_assert(sizeof(mdb_size_t) == sizeof(size_t), "mdb_size_t contract");

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

_Static_assert(MDB_PREV_MULTIPLE == 18, "MDB_cursor_op contract");
_Static_assert(MDB_APPENDDUP == 0x40000, "flag value contract");
_Static_assert(MDB_KEYEXIST == -30799, "errno range contract");
