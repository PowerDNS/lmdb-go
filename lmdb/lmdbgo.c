/* lmdbgo.c
 * Dispatch shims for github.com/PowerDNS/lmdb-go/v2/lmdb. See lmdbgo.h.
 * */
#include "lmdb.h"
#include "lmdbgo.h"
#include "_cgo_export.h"

/* Declare both engines' renamed symbol sets against the canonical types.
 * On Windows only the 0.9 engine exists (LMDB 1.0 is broken there and its
 * tree carries a //go:build !windows constraint). */
#define LMDBGO_SYM(name) mdb09_##name
#include "lmdbgo_dispatch.h"
#undef LMDBGO_SYM
#ifndef _WIN32
#define LMDBGO_SYM(name) mdb10_##name
#include "lmdbgo_dispatch.h"
#undef LMDBGO_SYM
#endif

/* Engine dispatch: one predicted branch to a direct call. ver is 9 or 10,
 * copied from the Go Env/Txn/Cursor. Valid for void calls too (a conditional
 * operator with two void operands is void). */
#ifdef _WIN32
# define LMDBGO_DISPATCH(ver, name, ...) mdb09_##name(__VA_ARGS__)
#else
# define LMDBGO_DISPATCH(ver, name, ...) \
    ((ver) >= 10 ? mdb10_##name(__VA_ARGS__) : mdb09_##name(__VA_ARGS__))
#endif

#define LMDBGO_SET_VAL(val, size, data) \
    *(val) = (MDB_val){.mv_size = (size), .mv_data = (data)}

/* Environment. */

int lmdbgo2_mdb_env_create(int ver, MDB_env **env) {
    return LMDBGO_DISPATCH(ver, env_create, env);
}

int lmdbgo2_mdb_env_open(int ver, MDB_env *env, const char *path, unsigned int flags, mdb_mode_t mode) {
    return LMDBGO_DISPATCH(ver, env_open, env, path, flags, mode);
}

void lmdbgo2_mdb_env_close(int ver, MDB_env *env) {
    LMDBGO_DISPATCH(ver, env_close, env);
}

int lmdbgo2_mdb_env_copy(int ver, MDB_env *env, const char *path) {
    return LMDBGO_DISPATCH(ver, env_copy, env, path);
}

int lmdbgo2_mdb_env_copy2(int ver, MDB_env *env, const char *path, unsigned int flags) {
    return LMDBGO_DISPATCH(ver, env_copy2, env, path, flags);
}

int lmdbgo2_mdb_env_copyfd(int ver, MDB_env *env, mdb_filehandle_t fd) {
    return LMDBGO_DISPATCH(ver, env_copyfd, env, fd);
}

int lmdbgo2_mdb_env_copyfd2(int ver, MDB_env *env, mdb_filehandle_t fd, unsigned int flags) {
    return LMDBGO_DISPATCH(ver, env_copyfd2, env, fd, flags);
}

int lmdbgo2_mdb_env_sync(int ver, MDB_env *env, int force) {
    return LMDBGO_DISPATCH(ver, env_sync, env, force);
}

int lmdbgo2_mdb_env_stat(int ver, MDB_env *env, MDB_stat *stat) {
    return LMDBGO_DISPATCH(ver, env_stat, env, stat);
}

int lmdbgo2_mdb_env_info(int ver, MDB_env *env, MDB_envinfo *info) {
    return LMDBGO_DISPATCH(ver, env_info, env, info);
}

int lmdbgo2_mdb_env_get_flags(int ver, MDB_env *env, unsigned int *flags) {
    return LMDBGO_DISPATCH(ver, env_get_flags, env, flags);
}

int lmdbgo2_mdb_env_set_flags(int ver, MDB_env *env, unsigned int flags, int onoff) {
    return LMDBGO_DISPATCH(ver, env_set_flags, env, flags, onoff);
}

int lmdbgo2_mdb_env_get_path(int ver, MDB_env *env, const char **path) {
    return LMDBGO_DISPATCH(ver, env_get_path, env, path);
}

int lmdbgo2_mdb_env_get_fd(int ver, MDB_env *env, mdb_filehandle_t *fd) {
    return LMDBGO_DISPATCH(ver, env_get_fd, env, fd);
}

int lmdbgo2_mdb_env_set_mapsize(int ver, MDB_env *env, size_t size) {
    return LMDBGO_DISPATCH(ver, env_set_mapsize, env, size);
}

int lmdbgo2_mdb_env_set_maxreaders(int ver, MDB_env *env, unsigned int readers) {
    return LMDBGO_DISPATCH(ver, env_set_maxreaders, env, readers);
}

int lmdbgo2_mdb_env_get_maxreaders(int ver, MDB_env *env, unsigned int *readers) {
    return LMDBGO_DISPATCH(ver, env_get_maxreaders, env, readers);
}

int lmdbgo2_mdb_env_set_maxdbs(int ver, MDB_env *env, MDB_dbi dbs) {
    return LMDBGO_DISPATCH(ver, env_set_maxdbs, env, dbs);
}

int lmdbgo2_mdb_env_get_maxkeysize(int ver, MDB_env *env) {
    return LMDBGO_DISPATCH(ver, env_get_maxkeysize, env);
}

int lmdbgo2_mdb_reader_check(int ver, MDB_env *env, int *dead) {
    return LMDBGO_DISPATCH(ver, reader_check, env, dead);
}

/* Transactions. */

int lmdbgo2_mdb_txn_begin(int ver, MDB_env *env, MDB_txn *parent, unsigned int flags, MDB_txn **txn) {
    return LMDBGO_DISPATCH(ver, txn_begin, env, parent, flags, txn);
}

int lmdbgo2_mdb_txn_commit(int ver, MDB_txn *txn) {
    return LMDBGO_DISPATCH(ver, txn_commit, txn);
}

void lmdbgo2_mdb_txn_abort(int ver, MDB_txn *txn) {
    LMDBGO_DISPATCH(ver, txn_abort, txn);
}

void lmdbgo2_mdb_txn_reset(int ver, MDB_txn *txn) {
    LMDBGO_DISPATCH(ver, txn_reset, txn);
}

int lmdbgo2_mdb_txn_renew(int ver, MDB_txn *txn) {
    return LMDBGO_DISPATCH(ver, txn_renew, txn);
}

size_t lmdbgo2_mdb_txn_id(int ver, MDB_txn *txn) {
    return LMDBGO_DISPATCH(ver, txn_id, txn);
}

/* Databases. */

int lmdbgo2_mdb_dbi_open(int ver, MDB_txn *txn, const char *name, unsigned int flags, MDB_dbi *dbi) {
    return LMDBGO_DISPATCH(ver, dbi_open, txn, name, flags, dbi);
}

int lmdbgo2_mdb_dbi_flags(int ver, MDB_txn *txn, MDB_dbi dbi, unsigned int *flags) {
    return LMDBGO_DISPATCH(ver, dbi_flags, txn, dbi, flags);
}

void lmdbgo2_mdb_dbi_close(int ver, MDB_env *env, MDB_dbi dbi) {
    LMDBGO_DISPATCH(ver, dbi_close, env, dbi);
}

int lmdbgo2_mdb_drop(int ver, MDB_txn *txn, MDB_dbi dbi, int del) {
    return LMDBGO_DISPATCH(ver, drop, txn, dbi, del);
}

int lmdbgo2_mdb_stat(int ver, MDB_txn *txn, MDB_dbi dbi, MDB_stat *stat) {
    return LMDBGO_DISPATCH(ver, stat, txn, dbi, stat);
}

/* Cursors. */

int lmdbgo2_mdb_cursor_open(int ver, MDB_txn *txn, MDB_dbi dbi, MDB_cursor **cursor) {
    return LMDBGO_DISPATCH(ver, cursor_open, txn, dbi, cursor);
}

int lmdbgo2_mdb_cursor_renew(int ver, MDB_txn *txn, MDB_cursor *cursor) {
    return LMDBGO_DISPATCH(ver, cursor_renew, txn, cursor);
}

void lmdbgo2_mdb_cursor_close(int ver, MDB_cursor *cursor) {
    LMDBGO_DISPATCH(ver, cursor_close, cursor);
}

int lmdbgo2_mdb_cursor_get(int ver, MDB_cursor *cursor, MDB_val *key, MDB_val *data, MDB_cursor_op op) {
    return LMDBGO_DISPATCH(ver, cursor_get, cursor, key, data, op);
}

int lmdbgo2_mdb_cursor_del(int ver, MDB_cursor *cursor, unsigned int flags) {
    return LMDBGO_DISPATCH(ver, cursor_del, cursor, flags);
}

int lmdbgo2_mdb_cursor_count(int ver, MDB_cursor *cursor, size_t *countp) {
    return LMDBGO_DISPATCH(ver, cursor_count, cursor, countp);
}

MDB_dbi lmdbgo2_mdb_cursor_dbi(int ver, MDB_cursor *cursor) {
    return LMDBGO_DISPATCH(ver, cursor_dbi, cursor);
}

/* Misc. */

char *lmdbgo2_mdb_strerror(int err) {
#ifdef _WIN32
    return mdb09_strerror(err);
#else
    /* The 1.0 table is a superset; shared codes have equivalent messages. */
    return mdb10_strerror(err);
#endif
}

char *lmdbgo2_mdb_version(int ver, int *major, int *minor, int *patch) {
    return LMDBGO_DISPATCH(ver, version, major, minor, patch);
}

/* Reader-list bridge. */

static int lmdbgo2_mdb_msg_func_proxy(const char *msg, void *ctx) {
    /* wrap msg and call the bridge function exported from msgfunc.go. */
    lmdbgo_ConstCString s;
    s.p = msg;
    return lmdbgo2MDBMsgFuncBridge(s, (size_t)ctx);
}

int lmdbgo2_mdb_reader_list(int ver, MDB_env *env, size_t ctx) {
    /* list readers using a static proxy function that does dynamic dispatch on
     * ctx. */
    if (ctx)
        return LMDBGO_DISPATCH(ver, reader_list, env, &lmdbgo2_mdb_msg_func_proxy, (void *)ctx);
    return LMDBGO_DISPATCH(ver, reader_list, env, 0, (void *)ctx);
}

/* Data-path shims (see lmdbgo.h for why these take char* instead of void*). */

int lmdbgo2_mdb_del(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn) {
    MDB_val key, val;
    LMDBGO_SET_VAL(&key, kn, kdata);
    LMDBGO_SET_VAL(&val, vn, vdata);
    return LMDBGO_DISPATCH(ver, del, txn, dbi, &key, &val);
}

int lmdbgo2_mdb_get(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, MDB_val *val) {
    MDB_val key;
    LMDBGO_SET_VAL(&key, kn, kdata);
    return LMDBGO_DISPATCH(ver, get, txn, dbi, &key, val);
}

int lmdbgo2_mdb_put2(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags) {
    MDB_val key, val;
    LMDBGO_SET_VAL(&key, kn, kdata);
    LMDBGO_SET_VAL(&val, vn, vdata);
    return LMDBGO_DISPATCH(ver, put, txn, dbi, &key, &val, flags);
}

int lmdbgo2_mdb_put1(int ver, MDB_txn *txn, MDB_dbi dbi, char *kdata, size_t kn, MDB_val *val, unsigned int flags) {
    MDB_val key;
    LMDBGO_SET_VAL(&key, kn, kdata);
    return LMDBGO_DISPATCH(ver, put, txn, dbi, &key, val, flags);
}

int lmdbgo2_mdb_cursor_put2(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags) {
    MDB_val key, val;
    LMDBGO_SET_VAL(&key, kn, kdata);
    LMDBGO_SET_VAL(&val, vn, vdata);
    return LMDBGO_DISPATCH(ver, cursor_put, cur, &key, &val, flags);
}

int lmdbgo2_mdb_cursor_put1(int ver, MDB_cursor *cur, char *kdata, size_t kn, MDB_val *val, unsigned int flags) {
    MDB_val key;
    LMDBGO_SET_VAL(&key, kn, kdata);
    return LMDBGO_DISPATCH(ver, cursor_put, cur, &key, val, flags);
}

int lmdbgo2_mdb_cursor_putmulti(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, size_t vstride, unsigned int flags) {
    MDB_val key, val[2];
    LMDBGO_SET_VAL(&key, kn, kdata);
    LMDBGO_SET_VAL(&(val[0]), vstride, vdata);
    LMDBGO_SET_VAL(&(val[1]), vn, 0);
    return LMDBGO_DISPATCH(ver, cursor_put, cur, &key, &val[0], flags);
}

int lmdbgo2_mdb_cursor_get1(int ver, MDB_cursor *cur, char *kdata, size_t kn, MDB_val *key, MDB_val *val, MDB_cursor_op op) {
    LMDBGO_SET_VAL(key, kn, kdata);
    return LMDBGO_DISPATCH(ver, cursor_get, cur, key, val, op);
}

int lmdbgo2_mdb_cursor_get2(int ver, MDB_cursor *cur, char *kdata, size_t kn, char *vdata, size_t vn, MDB_val *key, MDB_val *val, MDB_cursor_op op) {
    LMDBGO_SET_VAL(key, kn, kdata);
    LMDBGO_SET_VAL(val, vn, vdata);
    return LMDBGO_DISPATCH(ver, cursor_get, cur, key, val, op);
}
