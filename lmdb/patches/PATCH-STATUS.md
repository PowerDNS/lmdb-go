# Patch status

Per-patch record: what it fixes, upstream status, origin, and re-derivation
notes for LMDB version bumps. See [README.md](README.md) for how patches are
applied.

## Stream 09 (LMDB 0.9.x)

*(no patches yet)*

## Stream 10 (LMDB 1.0.x)

### cursor-close-after-txn.patch

- **Fixes:** use-after-free in `mdb_cursor_close` when a readonly cursor is
  closed after its transaction ended — a sequence both LMDB versions document
  as legal (see `mdb_cursor_open`/`mdb_cursor_close` docs) and that lmdb-go's
  cursor finalizers rely on. Pristine 1.0.0 evaluates the `MDB_CURSOR_UNREF`
  condition, which reads `mc_txn->mt_env->me_flags`, before the "read-only
  txn may have been freed already" guard. Reproducible: `TestCursor_Close_afterTxn`
  and `TestTxn_finalizer` SIGSEGV on the 1.0 engine without this patch, and
  `scripts/test-cursor-close-asan.sh` shows the heap-use-after-free under
  AddressSanitizer outside cgo.
- **Fix shape:** gate the unref on `C_UNTRACK`. Tracked (write-txn) cursors
  always have a live txn. For untracked (readonly) cursors the unref was a
  no-op unless `MDB_REMAP_CHUNKS` is set — which lmdb-go never sets — and
  with it, skipping only keeps a live readonly txn's page refs until
  reset/renew/end.
- **Upstream status:** not reported (OpenLDAP's contribution policy excludes
  AI-assisted contributions; see PowerDNS/lmdb-go#41 for context). Upstream's
  documented `MDB_RPAGE_CACHE=0` opt-out does not compile in pristine 1.0.0.
- **Re-derivation after a 1.0.x bump:** locate `mdb_cursor_close`, confirm
  whether upstream still unrefs before the tracked check; if fixed upstream,
  drop the patch. Verify with `scripts/test-cursor-close-asan.sh` and
  `LMDBGO_DEFAULT_FORMAT=10 go test ./lmdb`.
- **1.0.1 (2026-09-03):** still needed; the `mdb_cursor_close` region is
  unchanged upstream and the patch applies with a 25-line offset.
