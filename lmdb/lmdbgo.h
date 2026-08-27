/* lmdbgo.h
 * Dispatch shims for github.com/PowerDNS/lmdb-go/v2/lmdb. These functions
 * have no compatibility guarantees and may be modified or deleted without
 * warning.
 *
 * Two LMDB engines (streams) are linked into every build, their extern
 * symbols renamed to mdb09_* / mdb10_*. Go cannot call C function pointers,
 * so every LMDB call crosses cgo through one of the thin lmdbgo2_* shims
 * below, which take the engine selector `ver` (9 or 10, stored on the Go
 * Env/Txn/Cursor) as their first argument and branch to the prefixed symbol.
 *
 * The lmdbgo2_ prefix (rather than v1's lmdbgo_) is deliberate: lmdb-go v1
 * and v2 must be co-importable in one binary, and cgo does not hide C
 * symbols, so no v2 symbol may collide with v1's lmdbgo_* helpers or its
 * unprefixed mdb_* engine.
 *
 * All types below are the canonical surface from the shim lmdb.h (the LMDB
 * 1.0 header). The per-stream check TUs (lmdbgo_check09.c/lmdbgo_check10.c)
 * prove at compile time that both engines' real prototypes and struct
 * layouts match this surface.
 */
#ifndef LMDBGO2_H
#define LMDBGO2_H

#include "lmdb.h"

/* Environment. */
int lmdbgo2_mdb_env_create(int ver, MDB_env **env);
int lmdbgo2_mdb_env_open(int ver, MDB_env *env, const char *path, unsigned int flags, mdb_mode_t mode);
void lmdbgo2_mdb_env_close(int ver, MDB_env *env);
int lmdbgo2_mdb_env_copy(int ver, MDB_env *env, const char *path);
int lmdbgo2_mdb_env_copy2(int ver, MDB_env *env, const char *path, unsigned int flags);
int lmdbgo2_mdb_env_copyfd(int ver, MDB_env *env, mdb_filehandle_t fd);
int lmdbgo2_mdb_env_copyfd2(int ver, MDB_env *env, mdb_filehandle_t fd, unsigned int flags);
int lmdbgo2_mdb_env_sync(int ver, MDB_env *env, int force);
int lmdbgo2_mdb_env_stat(int ver, MDB_env *env, MDB_stat *stat);
int lmdbgo2_mdb_env_info(int ver, MDB_env *env, MDB_envinfo *info);
int lmdbgo2_mdb_env_get_flags(int ver, MDB_env *env, unsigned int *flags);
int lmdbgo2_mdb_env_set_flags(int ver, MDB_env *env, unsigned int flags, int onoff);
int lmdbgo2_mdb_env_get_path(int ver, MDB_env *env, const char **path);
int lmdbgo2_mdb_env_get_fd(int ver, MDB_env *env, mdb_filehandle_t *fd);
int lmdbgo2_mdb_env_set_mapsize(int ver, MDB_env *env, size_t size);
int lmdbgo2_mdb_env_set_maxreaders(int ver, MDB_env *env, unsigned int readers);
int lmdbgo2_mdb_env_get_maxreaders(int ver, MDB_env *env, unsigned int *readers);
int lmdbgo2_mdb_env_set_maxdbs(int ver, MDB_env *env, MDB_dbi dbs);
int lmdbgo2_mdb_env_get_maxkeysize(int ver, MDB_env *env);
int lmdbgo2_mdb_reader_check(int ver, MDB_env *env, int *dead);

/* Transactions. */
int lmdbgo2_mdb_txn_begin(int ver, MDB_env *env, MDB_txn *parent, unsigned int flags, MDB_txn **txn);
int lmdbgo2_mdb_txn_commit(int ver, MDB_txn *txn);
void lmdbgo2_mdb_txn_abort(int ver, MDB_txn *txn);
void lmdbgo2_mdb_txn_reset(int ver, MDB_txn *txn);
int lmdbgo2_mdb_txn_renew(int ver, MDB_txn *txn);
size_t lmdbgo2_mdb_txn_id(int ver, MDB_txn *txn);

/* Databases. */
int lmdbgo2_mdb_dbi_open(int ver, MDB_txn *txn, const char *name, unsigned int flags, MDB_dbi *dbi);
int lmdbgo2_mdb_dbi_flags(int ver, MDB_txn *txn, MDB_dbi dbi, unsigned int *flags);
void lmdbgo2_mdb_dbi_close(int ver, MDB_env *env, MDB_dbi dbi);
int lmdbgo2_mdb_drop(int ver, MDB_txn *txn, MDB_dbi dbi, int del);
int lmdbgo2_mdb_stat(int ver, MDB_txn *txn, MDB_dbi dbi, MDB_stat *stat);

/* Cursors. */
int lmdbgo2_mdb_cursor_open(int ver, MDB_txn *txn, MDB_dbi dbi, MDB_cursor **cursor);
int lmdbgo2_mdb_cursor_renew(int ver, MDB_txn *txn, MDB_cursor *cursor);
void lmdbgo2_mdb_cursor_close(int ver, MDB_cursor *cursor);
int lmdbgo2_mdb_cursor_get(int ver, MDB_cursor *cursor, MDB_val *key, MDB_val *data, MDB_cursor_op op);
int lmdbgo2_mdb_cursor_del(int ver, MDB_cursor *cursor, unsigned int flags);
int lmdbgo2_mdb_cursor_count(int ver, MDB_cursor *cursor, size_t *countp);
MDB_dbi lmdbgo2_mdb_cursor_dbi(int ver, MDB_cursor *cursor);

/* Misc. mdb_strerror takes no ver: the LMDB 1.0 message table is a superset
 * and messages for shared codes are equivalent. mdb_version reports the
 * engine selected by ver. */
char *lmdbgo2_mdb_strerror(int err);
char *lmdbgo2_mdb_version(int ver, int *major, int *minor, int *patch);

/* Proxy functions for lmdb get/put operations. The functions are defined to
 * take char* values instead of void* to keep cgo from checking their data for
 * nested pointers and causing a couple of allocations per argument.
 *
 * See these issues for more information about the problem and the decision:
 *      https://github.com/golang/go/issues/14387
 *      https://github.com/golang/go/issues/15048
 *      https://github.com/PowerDNS/lmdb-go/issues/63
 * */
int lmdbgo2_mdb_del(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn);
int lmdbgo2_mdb_get(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, MDB_val *val);
int lmdbgo2_mdb_put1(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, MDB_val *val, unsigned int flags);
int lmdbgo2_mdb_put2(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags);
int lmdbgo2_mdb_cursor_put1(int ver, MDB_cursor *cur, char *kdata, size_t kn, MDB_val *val, unsigned int flags);
int lmdbgo2_mdb_cursor_put2(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags);
int lmdbgo2_mdb_cursor_putmulti(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, size_t vstride, unsigned int flags);
int lmdbgo2_mdb_cursor_get1(int ver, MDB_cursor *cur, char *kdata, size_t kn, MDB_val *key, MDB_val *val, MDB_cursor_op op);
int lmdbgo2_mdb_cursor_get2(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, MDB_val *key, MDB_val *val, MDB_cursor_op op);

/* ConstCString wraps a null-terminated (const char *) because Go's type system
 * does not represent the 'const' qualifier directly on a function argument and
 * causes warnings to be emitted during linking.
 * */
typedef struct{ const char *p; } lmdbgo_ConstCString;

/* lmdbgo2_mdb_reader_list is a proxy for mdb_reader_list that uses a special
 * mdb_msg_func proxy function to relay messages over the
 * lmdbgo2MDBMsgFuncBridge external Go func.
 * */
int lmdbgo2_mdb_reader_list(int ver, MDB_env *env, size_t ctx);

#endif
