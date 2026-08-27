/* lmdbgo_dispatch.h
 *
 * Declares one vendored LMDB engine's renamed extern symbols against the
 * including translation unit's own types.
 *
 * NO include guard: this header is included once per engine prefix. The
 * consumer defines LMDBGO_SYM(name) to map a bare LMDB function name (without
 * the mdb_ prefix) to the engine symbol, includes this header, and #undefs
 * LMDBGO_SYM again:
 *
 *   - lmdbgo.c includes it twice (mdb09_* and mdb10_*), against the canonical
 *     LMDB 1.0 header types, to implement the lmdbgo2_* dispatch shims.
 *   - lmdbgo_check09.c / lmdbgo_check10.c include it once each, in a TU where
 *     the stream's own pristine header (with renames applied) has already
 *     declared the same symbols. A same-TU redeclaration with incompatible
 *     types is a compile error, so any prototype drift between an engine and
 *     the canonical surface is caught at build time (e.g. this proves LMDB
 *     1.0's mdb_size_t parameters are the same type as the canonical size_t
 *     declarations on this platform).
 *
 * The list covers exactly the LMDB surface the Go binding uses. When adding a
 * Go binding for a new LMDB function, add its declaration here (it is then
 * automatically prototype-checked against both engines).
 */

/* Environment. */
extern int LMDBGO_SYM(env_create)(MDB_env **env);
extern int LMDBGO_SYM(env_open)(MDB_env *env, const char *path, unsigned int flags, mdb_mode_t mode);
extern void LMDBGO_SYM(env_close)(MDB_env *env);
extern int LMDBGO_SYM(env_copy)(MDB_env *env, const char *path);
extern int LMDBGO_SYM(env_copy2)(MDB_env *env, const char *path, unsigned int flags);
extern int LMDBGO_SYM(env_copyfd)(MDB_env *env, mdb_filehandle_t fd);
extern int LMDBGO_SYM(env_copyfd2)(MDB_env *env, mdb_filehandle_t fd, unsigned int flags);
extern int LMDBGO_SYM(env_sync)(MDB_env *env, int force);
extern int LMDBGO_SYM(env_stat)(MDB_env *env, MDB_stat *stat);
extern int LMDBGO_SYM(env_info)(MDB_env *env, MDB_envinfo *stat);
extern int LMDBGO_SYM(env_get_flags)(MDB_env *env, unsigned int *flags);
extern int LMDBGO_SYM(env_set_flags)(MDB_env *env, unsigned int flags, int onoff);
extern int LMDBGO_SYM(env_get_path)(MDB_env *env, const char **path);
extern int LMDBGO_SYM(env_get_fd)(MDB_env *env, mdb_filehandle_t *fd);
extern int LMDBGO_SYM(env_set_mapsize)(MDB_env *env, size_t size);
extern int LMDBGO_SYM(env_set_maxreaders)(MDB_env *env, unsigned int readers);
extern int LMDBGO_SYM(env_get_maxreaders)(MDB_env *env, unsigned int *readers);
extern int LMDBGO_SYM(env_set_maxdbs)(MDB_env *env, MDB_dbi dbs);
extern int LMDBGO_SYM(env_get_maxkeysize)(MDB_env *env);
extern int LMDBGO_SYM(reader_check)(MDB_env *env, int *dead);
extern int LMDBGO_SYM(reader_list)(MDB_env *env, MDB_msg_func *func, void *ctx);

/* Transactions. */
extern int LMDBGO_SYM(txn_begin)(MDB_env *env, MDB_txn *parent, unsigned int flags, MDB_txn **txn);
extern int LMDBGO_SYM(txn_commit)(MDB_txn *txn);
extern void LMDBGO_SYM(txn_abort)(MDB_txn *txn);
extern void LMDBGO_SYM(txn_reset)(MDB_txn *txn);
extern int LMDBGO_SYM(txn_renew)(MDB_txn *txn);
extern size_t LMDBGO_SYM(txn_id)(MDB_txn *txn);

/* Databases and data. */
extern int LMDBGO_SYM(dbi_open)(MDB_txn *txn, const char *name, unsigned int flags, MDB_dbi *dbi);
extern int LMDBGO_SYM(dbi_flags)(MDB_txn *txn, MDB_dbi dbi, unsigned int *flags);
extern void LMDBGO_SYM(dbi_close)(MDB_env *env, MDB_dbi dbi);
extern int LMDBGO_SYM(drop)(MDB_txn *txn, MDB_dbi dbi, int del);
extern int LMDBGO_SYM(stat)(MDB_txn *txn, MDB_dbi dbi, MDB_stat *stat);
extern int LMDBGO_SYM(get)(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data);
extern int LMDBGO_SYM(put)(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data, unsigned int flags);
extern int LMDBGO_SYM(del)(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data);

/* Cursors. */
extern int LMDBGO_SYM(cursor_open)(MDB_txn *txn, MDB_dbi dbi, MDB_cursor **cursor);
extern int LMDBGO_SYM(cursor_renew)(MDB_txn *txn, MDB_cursor *cursor);
extern void LMDBGO_SYM(cursor_close)(MDB_cursor *cursor);
extern int LMDBGO_SYM(cursor_get)(MDB_cursor *cursor, MDB_val *key, MDB_val *data, MDB_cursor_op op);
extern int LMDBGO_SYM(cursor_put)(MDB_cursor *cursor, MDB_val *key, MDB_val *data, unsigned int flags);
extern int LMDBGO_SYM(cursor_del)(MDB_cursor *cursor, unsigned int flags);
extern int LMDBGO_SYM(cursor_count)(MDB_cursor *cursor, size_t *countp);
extern MDB_dbi LMDBGO_SYM(cursor_dbi)(MDB_cursor *cursor);

/* Misc. */
extern char *LMDBGO_SYM(strerror)(int err);
extern char *LMDBGO_SYM(version)(int *major, int *minor, int *patch);

#ifdef LMDBGO_DECLARE_V10_ONLY
/* LMDB 1.0-only functions, declared only for the mdb10_ prefix. Reserved for
 * the deferred 1.0 feature release; the lmdbgo2_* shims for these return
 * ENOTSUP when ver < 10. Uncomment as Go bindings are added, e.g.:
 *
 * extern int LMDBGO_SYM(txn_prepare)(MDB_txn *txn);
 * extern int LMDBGO_SYM(env_set_encrypt)(MDB_env *env, MDB_enc_func *func, const MDB_val *key, unsigned int size);
 */
#endif
