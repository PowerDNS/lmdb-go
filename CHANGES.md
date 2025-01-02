# Release Change Log

## v1.9.3 (2025-01-02)

## What's Changed

* Upgrade LMDB to 0.9.33 & use official OpenLDAP repo as source for updates by @wojas in https://github.com/PowerDNS/lmdb-go/pull/35
* macOS: fix unknown warning -Wno-stringop-overflow by @wojas in https://github.com/PowerDNS/lmdb-go/pull/38
* CI: minimum CI Go version 1.22 to fix error by @wojas in https://github.com/PowerDNS/lmdb-go/pull/37
* Update docs for v1.9.2 & document release process by @wojas in https://github.com/PowerDNS/lmdb-go/pull/30

Thanks to @igor35hh @shane-ns1 @fiatjaf for issue reports and suggestions.

LMDB C library changes:

        LMDB 0.9.33 Release (2024/05/21)
               ITS#9037 mdb_page_search: fix error code when DBI record is missing
               ITS#10198 For win32, stop passing ignored parameter
               ITS#10212 Fix meta page usage by read only tools

        LMDB 0.9.32 Release (2024/01/29)
               ITS#9378 - Add ability to replay log and replay log tool
               ITS#10095 - partial revert of ITS#9278. The patch was incorrect and introduced numerous race conditions.
               ITS#10125 - mdb_load: fix cursor reinit in Append mode
               ITS#10137 - Allow users to define MDB_IDL_LOGN

**Full Changelog**: https://github.com/PowerDNS/lmdb-go/compare/v1.9.2...v1.9.3

## v1.9.2 (2023-12-07)

## What's Changed
* Fix: allocate C memory for MDB_val in readonly Txn by @wojas in https://github.com/PowerDNS/lmdb-go/pull/29
* Disable memcpy() overflow warnings in C compilation by @shane-ns1 in https://github.com/PowerDNS/lmdb-go/pull/19
* Replace reflect.SliceHeader with unsafe.Slice by @shane-ns1 in https://github.com/PowerDNS/lmdb-go/pull/21
* Add names to fields used in tests by @shane-ns1 in https://github.com/PowerDNS/lmdb-go/pull/22
* Fix file mode in Open() statement by @shane-ns1 in https://github.com/PowerDNS/lmdb-go/pull/24
* Fix issues reported by staticcheck in tests by @shane-ns1 in https://github.com/PowerDNS/lmdb-go/pull/26

## New Contributors
* @shane-ns1 made their first contribution in https://github.com/PowerDNS/lmdb-go/pull/19
* @fiatjaf reported https://github.com/PowerDNS/lmdb-go/issues/28

**Full Changelog**: https://github.com/PowerDNS/lmdb-go/compare/v1.9.1...v1.9.2

## v1.9.1 (2023-09-14)

- Fix release numbers in the readme by @pieterlexis in https://github.com/PowerDNS/lmdb-go/pull/10
- Define IntegerKey and IntegerDup flags by @wojas in https://github.com/PowerDNS/lmdb-go/pull/12
- Upgrade LMDB C lib from 0.9.28 to 0.9.29 by @wojas in https://github.com/PowerDNS/lmdb-go/pull/13
- Upgrade LMDB to 0.9.31 by @wojas in https://github.com/PowerDNS/lmdb-go/pull/14
- CI: switch to Github Actions by @wojas in https://github.com/PowerDNS/lmdb-go/pull/15

LMDB C library changes:

```
LMDB 0.9.31 Release (2023/07/10)
	ITS#8447 - Fix cursor_put(MDB_CURRENT) on DUPSORT DB with different sized data

LMDB 0.9.30 Release (2023/02/08)
	ITS#9806 - LMDB page_split: key threshold depends on page size
	ITS#9916 - avoid gcc optimization bug on sparc64 linux
	ITS#9919 - Mark infrequently used functions as cold
	ITS#9723 - clear C_EOF on cursor with MDB_FIRST_DUP
	ITS#9030 - Use sys/cachectl.h rather than asm/cachectl.h on mips

LMDB 0.9.29 Release (2021/03/16)
	ITS#9461 refix ITS#9376
	ITS#9500 fix regression from ITS#8662
```


**Full Changelog**: https://github.com/PowerDNS/lmdb-go/compare/v1.9.0...v0.9.1

## v1.9.0 (2021-04-20)

First release of this PowerDNS/lmdb-go fork.

- Renamed module from `github.com/bmatsuo/lmdb-go` to `github.com/PowerDNS/lmdb-go` (#3, PR #7)
- Add `go.mod` and fix tests and Travis CI on recent Go versions (#5, PR #6)
- Fix: Cursor.Put would write "\x00" instead of an empty value (#1, PR #2)
- Remove experimental, never released `exp/lmdbpool` package (PR #9)
- lmdb: Update LMDB C library to version 0.9.28 (#4).

```
	LMDB 0.9.28 Release (2021/02/04)
		ITS#8662 add -a append option to mdb_load
	
	LMDB 0.9.27 Release (2020/10/26)
		ITS#9376 fix repeated DUPSORT cursor deletes
	
	LMDB 0.9.26 Release (2020/08/11)
		ITS#9278 fix robust mutex cleanup for FreeBSD
	
	LMDB 0.9.25 Release (2020/01/30)
		ITS#9068 fix mdb_dump/load backslashes in printable content
		ITS#9118 add MAP_NOSYNC for FreeBSD
		ITS#9155 free mt_spill_pgs in non-nested txn on end
	
	LMDB 0.9.24 Release (2019/07/24)
		ITS#8969 Tweak mdb_page_split
		ITS#8975 WIN32 fix writemap set_mapsize crash
		ITS#9007 Fix loose pages in WRITEMAP
	
	LMDB 0.9.23 Release (2018/12/19)
		ITS#8756 Fix loose pages in dirty list
		ITS#8831 Fix mdb_load flag init
		ITS#8844 Fix mdb_env_close in forked process
		Documentation
			ITS#8857 mdb_cursor_del doesn't invalidate cursor
			ITS#8908 GET_MULTIPLE etc don't change passed in key
	
	LMDB 0.9.22 Release (2018/03/22)
		Fix MDB_DUPSORT alignment bug (ITS#8819)
		Fix regression with new db from 0.9.19 (ITS#8760)
		Fix liblmdb to build on Solaris (ITS#8612)
		Fix delete behavior with DUPSORT DB (ITS#8622)
		Fix mdb_cursor_get/mdb_cursor_del behavior (ITS#8722)
	
	LMDB 0.9.21 Release (2017/06/01)
		Fix xcursor after cursor_del (ITS#8622)
	
	LMDB 0.9.20 (Withdrawn)
		Fix mdb_load with escaped plaintext (ITS#8558)
		Fix mdb_cursor_last / mdb_put interaction (ITS#8557)
```


Changes predating the PowerDNS fork (up to 2017):

- Fix unsafe threading behavior in benchmarks (bmatsuo/lmdb-go#101)
- Update transactions no longer allocate `MDB_val` objects (bmatsuo/lmdb-go#102)
- Txn.Renew no longer clears the Txn finalizer -- prevents resource leaks (bmatsuo/lmdb-go#104)
- Txn.Pooled field added so that the Txn finalizer may work better with
  sync.Pool (bmatsuo/lmdb-go#104 bmatsuo/lmdb-go#105)
- Fixed a race in the Txn finalizer that could lead to a segfault (bmatsuo/lmdb-go#105)
- Txn.RunOp method added so that it is possible for other packages to create
  other flavors of managed transactions from scratch (bmatsuo/lmdb-go#105)
- Experimental package lmdbpool was added to make integration of lmdb and
  sync.Pool easier (bmatsuo/lmdb-go#104 bmatsuo/lmdb-go#105)

```
go get github.com/PowerDNS/lmdb-go/exp/lmdbpool
```

- Silence aggressive struct initializer warning from clang (bmatsuo/lmdb-go#107)
- Improved documentation regarding long-running transactions and dead readers
  (bmatsuo/lmdb-go#111)

## v1.8.0 (2017-02-10)

- lmdbscan: The package was moved out of the exp/ subtree and can now be
  considered stable and suitable for general use.
- lmdb: Update LMDB C library to version 0.9.19 (bmatsuo/lmdb-go#92).

```
	Fix mdb_env_cwalk cursor init (ITS#8424)
	Fix robust mutexes on Solaris 10/11 (ITS#8339)
	Tweak Win32 error message buffer
	Fix MDB_GET_BOTH on non-dup record (ITS#8393)
	Optimize mdb_drop
	Fix xcursors after mdb_cursor_del (ITS#8406)
	Fix MDB_NEXT_DUP after mdb_cursor_del (ITS#8412)
	Fix mdb_cursor_put resetting C_EOF (ITS#8489)
	Fix mdb_env_copyfd2 to return EPIPE on SIGPIPE (ITS#8504)
	Fix mdb_env_copy with empty DB (ITS#8209)
	Fix behaviors with fork (ITS#8505)
	Fix mdb_dbi_open with mainDB cursors (ITS#8542)
	Fix robust mutexes on kFreeBSD (ITS#8554)
	Fix utf8_to_utf16 error checks (ITS#7992)
	Fix F_NOCACHE on MacOS, error is non-fatal (ITS#7682)
	Build
		Make shared lib suffix overridable (ITS#8481)
	Documentation
		Cleanup doxygen nits
		Note reserved vs actual mem/disk usage
```

- lmdb: Fix resource leak in cursor tests (bcf4e9f).
- lmdb: Fix panic in Cursor.Get when using the Set op (bmatsuo/lmdb-go#96).
- docs: Improve documentation about when runtime.LockOSThread is required

## v1.7.0

- lmdb: Removed unnecessary import of the "math" package (bmatsuo/lmdb-go#70).
- lmdb: Removed direct dependency on the "fmt" package and reduced error
  related allocation (bmatsuo/lmdb-go#73).
- cmd/lmdb_stat: Fix transaction ID decoding and match output of `mdb_stat`
  1-to-1 (bmatsuo/lmdb-go#78).
- lmdb: fix compilation for 32-bit architectures (bmatsuo/lmdb-go#83).

## v1.6.0 (2016-04-07)

- lmdb: method Txn.ID() exposing mdb_txn_id. (bmatsuo/lmdb-go#47)
- lmdb: Env.ReaderList() returns an error if passed a nil function. (bmatsuo/lmdb-go#48)
- lmdbsync: realistic test of resizing functionality (bmatsuo/lmdb-go#7)
- lmdbsync: use context.Context instead of a hand-rolled Bag (bmatsuo/lmdb-go#51)
- lmdbsync: Handler Env is now an argument instead of a context value (bmatsuo/lmdb-go#52)
- lmdbsync: Changes to MapResizedHandler and its default values (bmatsuo/lmdb-go#54)
- lmdb: Fix CGO argument check panic for certain []byte values produced from a
  bytes.Buffer (bmatsuo/lmdb-go#56)
- lmdb: Support building the C library with support for the pwritev(2) system
  call (bmatsuo/lmdb-go#58)
- lmdb: Reuse MDB_val values within transactions to reduce allocations in
  transactions issuing multiple Get operations (bmatsuo/lmdb-go#61).
- lmdb: Avoid allocation and linear scan overhead on the cgo boundary for
  transaction operations (Get/Put and variants) (bmatsuo/lmdb-go#63).
- lmdb: Use a more portable internal conversion from C pointers to slices
  (bmatsuo/lmdb-go#67).

## v1.5.0

- lmdb: fix crash from bad interaction with Txn finalizer and Txn.Reset/.Renew.
- lmdb: Update the LMDB C library to 0.9.18

```
    Fix robust mutex detection on glibc 2.10-11 (ITS#8330)
    Fix page_search_root assert on FreeDB (ITS#8336)
    Fix MDB_APPENDDUP vs. rewrite(single item) (ITS#8334)
    Fix mdb_copy of large files on Windows
    Fix subcursor move after delete (ITS#8355)
    Fix mdb_midl_shirnk off-by-one (ITS#8363)
    Check for utf8_to_utf16 failures (ITS#7992)
    Catch strdup failure in mdb_dbi_open
    Build
        Additional makefile var tweaks (ITS#8169)
    Documentation
        Add Getting Started page
        Update WRITEMAP description
```

## v1.4.0

- development: The LMDB C library can be cloned under /lmdb -- it will be
  ignored.
- lmdb: Pass CFLAGS -Wno-format-extra-args to silence compilation warning (OS
  X).
- lmdb: Update the LMDB C library to 0.9.17

```
    Fix ITS#7377 catch calloc failure
    Fix ITS#8237 regression from ITS#7589
    Fix ITS#8238 page_split for DUPFIXED pages
    Fix ITS#8221 MDB_PAGE_FULL on delete/rebalance
    Fix ITS#8258 rebalance/split assert
    Fix ITS#8263 cursor_put cursor tracking
    Fix ITS#8264 cursor_del cursor tracking
    Fix ITS#8310 cursor_del cursor tracking
    Fix ITS#8299 mdb_del cursor tracking
    Fix ITS#8300 mdb_del cursor tracking
    Fix ITS#8304 mdb_del cursor tracking
    Fix ITS#7771 fakepage cursor tracking
    Fix ITS#7789 ensure mapsize >= pages in use
    Fix ITS#7971 mdb_txn_renew0() new reader slots
    Fix ITS#7969 use __sync_synchronize on non-x86
    Fix ITS#8311 page_split from update_key
    Fix ITS#8312 loose pages in nested txn
    Fix ITS#8313 mdb_rebalance dummy cursor
    Fix ITS#8315 dirty_room in nested txn
    Fix ITS#8323 dirty_list in nested txn
    Fix ITS#8316 page_merge cursor tracking
    Fix ITS#8321 cursor tracking
    Fix ITS#8319 mdb_load error messages
    Fix ITS#8320 mdb_load plaintext input
    Added mdb_txn_id() (ITS#7994)
    Added robust mutex support
    Miscellaneous cleanup/simplification
    Build
        Create install dirs if needed (ITS#8256)
        Fix ThreadProc decl on Win32/MSVC (ITS#8270)
        Added ssize_t typedef for MSVC (ITS#8067)
        Use ANSI apis on Windows (ITS#8069)
        Use O_SYNC if O_DSYNC,MDB_DSYNC are not defined (ITS#7209)
        Allow passing AR to make (ITS#8168)
        Allow passing mandir to make install (ITS#8169)
```


## v1.3.0

- all: Builds on Windows with passing tests. Fixes bmatsuo/lmdb-go#33.
- lmdb: Cursor.DBI returns "invalid" DBI if the cursor is closed. Fixes bmatsuo/lmdb-go#31.
- lmdb: Finalizers to prevent resource leaks. Fixes bmatsuo/lmdb-go#20.
- all: Internal test package for setting up, populating, and tearing down environments.
- lmdbscan: Fix panic in Scanner.Scan after Txn.OpenCursor fails. Fixes bmatsuo/lmdb-go#21.
- lmdbscan: Scanner.Set[Next] methods move the cursor and make the next
  Scanner.Scan a noop.  The changes should be backwards compatible. Fixes bmatsuo/lmdb-go#17.
- lmdb: Cgo calling convention meets rules set forth for go1.6. Fixes bmatsuo/lmdb-go#10.
- lmdb: add a "Package" code example that shows a complete workflow

## v1.2.0

- Many example tests replaced with simpler code examples.
- Lots of documentation fixes
- internal/lmdbcmd: simplify version printing
- lmdbscan: add method Scanner.Cursor() to deprecate Scanner.Del()
- lmdbscan: add tests for Scanner.Set and Scanner.SetNext
- lmdb: Implement Env.FD() method returning an open file descriptor
- lmdbgo.c: remove unnecessary `#include <string.h>`

## v1.1.1

- Lots of code examples.
- Lots of tests.
- Travis-CI enforcing code style using [`golint`](https://github.com/golang/lint)
- exp/lmdbscan: removed the scanner.Func type because it was unnecessary bloat.
- exp/lmdbsync: Tweak lmdbsync.HandlerChain semantics
- exp/lmdbsync: Rename type RetryTxn to ErrTxnRetry
- Move exp/cmd/lmdb_stat to path cmd/lmdb_stat because its purpose is know and
  it is essentially complete.
- Move exp/cmd/lmdb_copy to path cmd/lmdb_stat because its purpose is know and
  it is essentially complete.
- Add method Env.ReaderList using C shim.
- exp/lmdbsync: Simplified interface and behavior after tests.
- exp/lmdbsync: No longer restrict implementations of lmdbsync.Bag with an
  unexported method.
- exp/lmdbsync: Do not let users call Env.BeginTxn because it is
  unsynchronized.
- lmdb: methods Env.CopyFD and Env.CopyFDFlags
- lmdb: clean up Multi.Vals by using Multi.Val internally
- exp/lmdbsync: clean up lmdbsync.MapFullHandler and lmdbsync.MapResizedHandler
  godoc.
- exp/lmdbsync: document possible deadlocks with MapFullHandler and MapResizedHandler
- exp/cmp/lmdb_example_resize: simple program that auto-resizes a database
- exp/lmdbsync: fix infinite loop
- README.md: link fixes
