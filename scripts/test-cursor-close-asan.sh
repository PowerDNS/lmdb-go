#!/usr/bin/env bash
#
# Verify the cursor-close-after-txn patch on the vendored LMDB 1.0 tree with
# AddressSanitizer, outside cgo (Go's runtime does not compose with ASan).
#
# Builds the vendored (patched) 1.0 engine standalone plus a minimal
# reproducer of the documented-legal sequence: open a readonly txn, open a
# cursor, end the txn, then close the cursor. On pristine LMDB 1.0.0 this is
# a heap-use-after-free (mdb_cursor_close reads the freed txn while
# evaluating the reader page cache state); with the patch ASan must be quiet.

set -euo pipefail
cd "$(dirname "$0")/.."

tmp="$(mktemp -d "${TMPDIR:-/tmp}/lmdbgo-asan.XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

cat > "$tmp/repro.c" <<'EOF'
#include "rename_lmdb10.h"
#include "lmdb.h"
#include <stdio.h>
#include <stdlib.h>

#define CHECK(expr) do { int rc_ = (expr); if (rc_ != MDB_SUCCESS) { \
    fprintf(stderr, "%s: %s\n", #expr, mdb_strerror(rc_)); exit(1); } } while (0)

int main(int argc, char **argv)
{
    MDB_env *env;
    MDB_txn *txn;
    MDB_dbi dbi;
    MDB_cursor *cur;

    if (argc < 2) { fprintf(stderr, "usage: repro <dbdir>\n"); return 2; }
    CHECK(mdb_env_create(&env));
    CHECK(mdb_env_open(env, argv[1], MDB_NOTLS, 0644));
    CHECK(mdb_txn_begin(env, NULL, MDB_RDONLY, &txn));
    CHECK(mdb_dbi_open(txn, NULL, 0, &dbi));
    CHECK(mdb_cursor_open(txn, dbi, &cur));
    mdb_txn_abort(txn);
    /* Documented-legal: close the readonly cursor after its txn ended. */
    mdb_cursor_close(cur);
    mdb_env_close(env);
    puts("OK: cursor close after txn end is safe");
    return 0;
}
EOF

cc -fsanitize=address -g -O1 -w -I lmdb \
    lmdb/mdb_lmdb10.c lmdb/midl_lmdb10.c "$tmp/repro.c" -o "$tmp/repro"

mkdir "$tmp/db"
"$tmp/repro" "$tmp/db"
