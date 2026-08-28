#!/usr/bin/env bash
#
# Update one of the two vendored LMDB streams in lmdb/.
# Releases: https://git.openldap.org/openldap/openldap/-/tags?sort=updated_desc&search=LMDB_
#
# lmdb-go v2 vendors two LMDB versions side by side in lmdb/, BOTH compiled
# into every build and selected per environment at runtime:
#
#   stream 09 (LMDB 0.9.x): mdb_lmdb09.c midl_lmdb09.c lmdb_lmdb09.h midl_lmdb09.h
#   stream 10 (LMDB 1.0.x): mdb_lmdb10.c midl_lmdb10.c lmdb_lmdb10.h midl_lmdb10.h
#
# Each stream's extern symbols are renamed (mdb_* -> mdb09_*/mdb10_*) by a
# generated rename_lmdbNN.h so the streams can coexist in one binary (and
# alongside lmdb-go v1). The lmdb.h and midl.h files in lmdb/ are hand-written
# shims that dispatch between the streams; this script does not touch them.
#
# The vendored .c files are byte-identical upstream sources except for:
#   1. mechanically prepended lines:
#        stream 10 only:  //go:build !windows   (LMDB 1.0 is broken on Windows)
#        both streams:    #include "rename_lmdbNN.h"
#   2. the patch series from lmdb/patches/lmdbNN/*.patch, applied in order
#      (see lmdb/patches/README.md and PATCH-STATUS.md).
#
# After vendoring, the rename header is regenerated from nm output so symbol
# renames stay complete across LMDB updates.

set -euo pipefail
cd "$(dirname "$0")"

stream="${1:-}"
version="${2:-}"

usage() {
    echo "USAGE: $0 <stream> <desired-version>"
    echo "  stream: 09 (LMDB 0.9.x) or 10 (LMDB 1.0.x)"
    echo "  e.g.: $0 09 0.9.35   or   $0 10 1.0.0"
    echo "Check https://git.openldap.org/openldap/openldap/-/tags?sort=updated_desc&search=LMDB_ for available versions"
}

case "$stream" in
    09|10) ;;
    *) usage; exit 1 ;;
esac
[ -n "$version" ] || { usage; exit 1; }

hdr="lmdb/lmdb_lmdb${stream}.h"
get_define() { grep "^#define $1" "$hdr" | head -1 | awk '{print $3}'; }
get_version() { echo "$(get_define MDB_VERSION_MAJOR).$(get_define MDB_VERSION_MINOR).$(get_define MDB_VERSION_PATCH)"; }

cur_version="(none)"
[ -f "$hdr" ] && cur_version="$(get_version)"
echo "Current LMDB version in stream $stream: $cur_version"
echo

tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/lmdb-update.XXXXXX")"
trap 'rm -rf "$tmp_dir"' EXIT
echo "Temp dir: $tmp_dir"

curl -fL "https://git.openldap.org/openldap/openldap/-/archive/LMDB_${version}/openldap-LMDB_${version}.tar.gz" \
    | tar -C "$tmp_dir" -xz
src="$tmp_dir/openldap-LMDB_${version}/libraries/liblmdb"
[ -f "$src/mdb.c" ] || { echo "ERROR: mdb.c not found in download" >&2; exit 1; }

# Write the vendored sources: mechanically prepended lines + pristine source.
vendor_c() {
    local from="$1" to="$2"
    : > "$to"
    if [ "$stream" = "10" ]; then
        # LMDB 1.0 is broken on Windows upstream; the 0.9 engine is the only
        # one available there. cgo honors build constraints in .c files.
        printf '//go:build !windows\n\n' >> "$to"
    fi
    printf '#include "rename_lmdb%s.h"\n\n' "$stream" >> "$to"
    cat "$from" >> "$to"
}
vendor_c "$src/mdb.c"  "lmdb/mdb_lmdb${stream}.c"
vendor_c "$src/midl.c" "lmdb/midl_lmdb${stream}.c"
cp "$src/lmdb.h" "lmdb/lmdb_lmdb${stream}.h"
cp "$src/midl.h" "lmdb/midl_lmdb${stream}.h"
cp "$src/CHANGES" "CHANGES.lmdb${stream}.txt"

# Apply this stream's patch series, in order. Patches are diffs against the
# vendored files (i.e. including the prepended lines above).
shopt -s nullglob
patches=(lmdb/patches/lmdb${stream}/*.patch)
shopt -u nullglob
if [ "${#patches[@]}" -gt 0 ]; then
    echo
    echo "Applying ${#patches[@]} patch(es) for stream $stream:"
    for p in "${patches[@]}"; do
        echo "  $p"
        git apply --verbose "$p"
    done
else
    echo "No patches for stream $stream."
fi

# Regenerate the rename header from the freshly vendored (and patched)
# sources, so new/removed externs are picked up.
scripts/gen-rename.sh "$stream"

# Sanity: verify the prepended lines survived.
if [ "$stream" = "10" ]; then
    head -1 "lmdb/mdb_lmdb10.c" | grep -q '^//go:build !windows$' \
        || { echo "ERROR: //go:build constraint missing from mdb_lmdb10.c" >&2; exit 1; }
fi
for f in "lmdb/mdb_lmdb${stream}.c" "lmdb/midl_lmdb${stream}.c"; do
    head -4 "$f" | grep -q "^#include \"rename_lmdb${stream}.h\"\$" \
        || { echo "ERROR: rename include missing from $f" >&2; exit 1; }
done

echo
new_version="$(get_version)"
echo "New LMDB version in stream $stream: $new_version"
echo
echo "NOTE: - Include the upstream changelog from $cur_version to $new_version from"
echo "        CHANGES.lmdb${stream}.txt in our CHANGES.md."
echo "      - Run: scripts/check-defines.sh (shared-define drift between streams)"
echo "      - Run the full test suite for BOTH engines, see Makefile."
echo "      - If patches failed to apply, re-derive them and update"
echo "        lmdb/patches/PATCH-STATUS.md."
