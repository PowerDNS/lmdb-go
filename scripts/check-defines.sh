#!/usr/bin/env bash
#
# Guard against constant/flag drift between the two vendored LMDB headers.
#
# The Go bindings and the C dispatch layer use a single canonical header (the
# LMDB 1.0 surface) for both engines, which is only sound while every MDB_*
# define shared by the two headers has the same value. Names expected to
# differ are whitelisted (version identity and the error-code high-water
# mark). Run after every stream update; wired into CI.

set -euo pipefail
cd "$(dirname "$0")/.."

tmp="$(mktemp -d -t lmdbgo-defines)"
trap 'rm -rf "$tmp"' EXIT

for s in 09 10; do
    # Post-preprocessor truth for this platform; only MDB_* names matter
    # (system headers do not define MDB_*).
    cc -E -dM -x c "lmdb/lmdb_lmdb${s}.h" \
        | awk '/^#define MDB_/ { name=$2; sub(/^#define [^ ]+ ?/, ""); print name "\t" $0 }' \
        | sort > "$tmp/$s.defs"
done

# Names allowed to differ between streams.
whitelist='^MDB_VERSION|^MDB_LAST_ERRCODE$'

mismatch="$(join -t "$(printf '\t')" "$tmp/09.defs" "$tmp/10.defs" \
    | awk -F '\t' -v wl="$whitelist" '$1 !~ wl && $2 != $3 { print $1 ": 09=[" $2 "] 10=[" $3 "]" }')"

if [ -n "$mismatch" ]; then
    echo "ERROR: shared MDB_* defines differ between lmdb_lmdb09.h and lmdb_lmdb10.h:" >&2
    printf '%s\n' "$mismatch" >&2
    exit 1
fi

n09="$(wc -l < "$tmp/09.defs" | tr -d ' ')"
n10="$(wc -l < "$tmp/10.defs" | tr -d ' ')"
echo "OK: shared MDB_* defines match (09: $n09 defines, 10: $n10 defines)"
