# Local patches on the vendored LMDB sources

Each subdirectory holds the patch series for one vendored LMDB stream:

- `lmdb09/` — patches on the LMDB 0.9.x tree (`lmdb/*_lmdb09.*`)
- `lmdb10/` — patches on the LMDB 1.0.x tree (`lmdb/*_lmdb10.*`)

Patches are unified diffs against the **vendored** files (including the
mechanically prepended `//go:build` / `#include "rename_lmdbNN.h"` lines) and
are applied in lexical order by `update-lmdb.sh` with `git apply` after every
stream update. Because the patches are applied at vendoring time in this
repository, consumers of the Go module always receive patched sources —
`lmdb/patches/` itself is not part of the build.

Every patch must have an entry in [PATCH-STATUS.md](PATCH-STATUS.md)
describing what it fixes, why it is not (yet) upstream, and how to re-derive
it after an LMDB version bump.

## Attribution

Some patches in this collection are imported from
[py-lmdb](https://github.com/jnwatson/py-lmdb) (jnwatson/py-lmdb), which
maintains a valuable series of LMDB hardening/security patches (see
jnwatson/py-lmdb#472 and the `lib*/py-lmdb/` directories in that repository).
Imported patches keep a reference to their py-lmdb origin in their header and
in PATCH-STATUS.md. py-lmdb and its patches are distributed under the
OpenLDAP Public License, the same license as the LMDB sources they modify
(see `LICENSE.mdb.md`).
