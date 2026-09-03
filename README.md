# lmdb-go [![releases](https://img.shields.io/badge/release-v2-375eab.svg)](https://github.com/PowerDNS/lmdb-go/releases) [![C/v0.9.35+v1.0.1](https://img.shields.io/badge/C-v0.9.35%20%2B%20v1.0.1-555555.svg)](https://git.openldap.org/openldap/openldap/-/tags?search=LMDB_) [![Build Status](https://github.com/PowerDNS/lmdb-go/actions/workflows/go.yml/badge.svg?branch=master)]

> **Looking for lmdb-go v1?** The v1 series (`github.com/PowerDNS/lmdb-go`,
> LMDB 0.9 only) lives on the [`v1` branch](https://github.com/PowerDNS/lmdb-go/tree/v1)
> and keeps receiving `v1.x.y` releases. This branch is **v2**
> (`github.com/PowerDNS/lmdb-go/v2`).

Go bindings to the OpenLDAP Lightning Memory-Mapped Database (LMDB).

## What is v2?

lmdb-go v2 ships **both LMDB 0.9 and LMDB 1.0** in one Go module, and lets a
program use either of them, or both, at runtime.

LMDB 1.0 (released June 2026) stores data in a new on-disk format that LMDB
0.9 cannot read, and vice versa. Both C libraries use the same function
names, so a normal program can link only one of them: it either reads old
databases or new ones. v2 removes that choice by bundling both engines under
distinct symbol names and picking an engine **per environment** when it is
opened:

- **Opening an existing database needs no configuration.** v2 looks at the
  database file, detects whether it was written by LMDB 0.9 or 1.0, and uses
  the matching engine. This is done without opening the environment, because
  a failed open with the wrong LMDB version can damage the lock file.
- **New databases are created in the 0.9 format by default**, so they stay
  readable by software that still uses LMDB 0.9: distribution packages, the
  `mdb_dump`/`mdb_stat` command line tools, bindings for other languages.
  You can opt into the 1.0 format per environment, for the whole process, or
  through an environment variable (see below).
- **Both formats can be open at the same time in one process.** lmdb-go never
  converts a database in place. To migrate, open the old database with one
  `Env`, create the new one with the 1.0 format, and copy the data over.

### Creating a database in the LMDB 1.0 format

Request the 1.0 format on the `Env` before calling `Open`:

```go
env, err := lmdb.NewEnv()
if err != nil {
	// ...
}
defer env.Close()

// SetFormat only affects databases that Open creates. An existing database
// is always opened with the engine matching its on-disk format, and if that
// format is not the one requested here, Open fails with
// lmdb.ErrFormatConflict.
if err := env.SetFormat(lmdb.V10); err != nil {
	// ...
}
if err := env.SetMapSize(1 << 30); err != nil {
	// ...
}
if err := env.Open("/path/to/db", 0, 0644); err != nil {
	// ...
}

log.Printf("database uses the LMDB %s format", env.Format()) // "0.9" or "1.0"
```

To make 1.0 the default for every new database in the process, call
`lmdb.SetDefaultFormat(lmdb.V10)` at startup, or set
`LMDBGO_DEFAULT_FORMAT=10` in the process environment. A call to
`SetDefaultFormat` wins over the environment variable, and `SetFormat` on an
`Env` wins over both. `lmdb.SniffFormat` reports the format of a database
file without opening it.

## Upgrading from v1

The Go API of v2 is drop-in compatible with v1. For most projects, upgrading
is only an import path change:

1. Replace `github.com/PowerDNS/lmdb-go/` with `github.com/PowerDNS/lmdb-go/v2/`
   in your imports, for every package you use (`lmdb`, `lmdbscan`,
   `exp/lmdbsync`). With GNU sed:

   ```sh
   find . -name '*.go' -exec sed -i 's#github.com/PowerDNS/lmdb-go/#github.com/PowerDNS/lmdb-go/v2/#g' {} +
   ```

   (On macOS use `sed -i ''`.)

2. Update `go.mod`:

   ```sh
   go get github.com/PowerDNS/lmdb-go/v2@latest
   go mod tidy
   ```

Nothing changes on disk. Your existing databases are in the 0.9 format and
v2 keeps opening them with the 0.9 engine, and new databases stay in the 0.9
format unless you opt into 1.0 as shown above.

Things to be aware of:

- **Go 1.21 or newer** is required.
- **v1 and v2 can be linked into the same binary.** Dependencies that still
  import v1 keep working next to your v2 code, so not everything has to move
  at once. (An `Env` cannot be shared between the two, of course.)
- **A few edge cases now return errors instead of panicking or crashing.**
  Methods on a closed cursor return an `EINVAL` error, using an `Env` before
  `Open` (for example `Stat` or `BeginTxn`) returns a clean error, and a
  failed `Open` can be retried.
- **`lmdb.Version` now reports the newest bundled engine** (1.0.x), with a
  release string naming both. Use `env.EngineVersion()` for the engine driving
  a specific environment. `env.MaxKeySize()` also depends on the engine: 511
  bytes on 0.9, and derived from the page size on 1.0 (8122 for 4K pages).
- **Windows builds contain only the 0.9 engine**, because upstream LMDB 1.0
  is currently broken there. Requesting the 1.0 format fails with a clear
  error.
- The `lmdb_stat` and `lmdb_copy` commands in [cmd/](cmd/) work with
  databases of either format.

The complete list of differences is in [CHANGES.md](CHANGES.md).

## Technical details

- **Bundled versions.** v2 currently bundles LMDB 0.9.35 and LMDB 1.0.1. As
  in v1, the bundled LMDB sources are always statically linked; dynamic
  linking against a system liblmdb is not supported.
- **Symbol renaming and dispatch.** The two vendored trees are compiled behind
  generated symbol-rename headers (`mdb09_*` / `mdb10_*`), and every LMDB
  call goes through a thin C shim that dispatches on a per-environment engine
  selector. The overhead versus v1 is one `int` argument plus a predicted
  branch per call: on the same 0.9 engine, the benchmark suite runs within
  noise of v1 and allocations are byte-identical. All v2 C symbols are
  namespaced (`lmdbgo2_*`), which is what makes linking v1 and v2 into one
  binary possible; `tests/coexist/` verifies this in CI.
- **Engine performance.** The engines themselves differ: on the 1.0 engine,
  creating a read-only transaction is noticeably more expensive than on 0.9
  (reader page cache setup), while renewed or pooled read-only transactions
  and write paths stay within a few percent, and sub-transactions are faster.
  Read-heavy workloads on the 1.0 engine benefit from `Txn.Renew` or pooling,
  as the package documentation already recommends.
- **Format detection.** `SniffFormat` reads the data file's meta page and
  distinguishes the formats by their data version number (1 for 0.9, 3 for
  1.0). Files in the lmdb-js prerelease format (data version 2) are rejected
  with `ErrFormatUnsupported`. A missing or empty data file means the
  environment does not exist yet and `Open` will create it.
- **Runtime introspection.** `Env.Format` reports the format (and thus the
  engine) driving an environment; `Env.EngineVersion` reports the exact
  upstream version of that engine.
- **LMDB 1.0 additions.** The LMDB 1.0 error constants (`Problem`,
  `BadChecksum`, `TxnPending`, ...) are always defined; the 0.9 engine never
  returns them. The `PrevSnapshot` open flag (`MDB_PREVSNAPSHOT`) opens the
  environment at its previous snapshot, which can help recover from some
  kinds of corruption; the 0.9 engine rejects it with `EINVAL`.
- **Local patches.** The vendored trees can carry local patches (see
  [lmdb/patches/](lmdb/patches/)), applied by `update-lmdb.sh` when a stream
  is updated. v2 ships a fix for a use-after-free in LMDB 1.0.x when a
  read-only cursor is closed after its transaction has ended, a sequence that
  LMDB documents as legal and that lmdb-go finalizers rely on.

## About this fork

This fork was created after the upstream repository [bmatsuo/lmdb-go](https://github.com/bmatsuo/lmdb-go)
went without updates for 4 years. We would like to express our thanks to
*bmatsuo* for creating and maintaining this repository until 2017.

We decided to rename this repository to allow usage without `replace`
directives in the `go.mod`, as we do not expect this to be a temporary
fork. This also makes it clear that new versions released here are
different from any upstream versions.

To use this package, update all your import paths from
`github.com/bmatsuo/lmdb-go` to `github.com/PowerDNS/lmdb-go` (v1) or
`github.com/PowerDNS/lmdb-go/v2` (v2). This affects all versions starting
from 1.9.0.

Note that the experimental, never released `exp/lmdbpool` package has
been removed in this fork.

## Packages

Functionality is logically divided into several packages.  Applications will
usually need to import **lmdb** but may import other packages on an as needed
basis.

Packages in the `exp/` directory are not stable and may change without warning.
That said, they are generally usable if application dependencies are managed
and pinned by tag/commit.

Developers concerned with package stability should consult the documentation.

#### lmdb [![GoDoc](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdb?status.svg)](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdb) [![stable](https://img.shields.io/badge/stability-stable-brightgreen.svg)](#user-content-versioning-and-stability)

```go
import "github.com/PowerDNS/lmdb-go/v2/lmdb"
```

Core bindings allowing low-level access to LMDB.

#### lmdbscan [![GoDoc](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdbscan?status.svg)](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdbscan) [![stable](https://img.shields.io/badge/stability-stable-brightgreen.svg)](#user-content-versioning-and-stability)

```go
import "github.com/PowerDNS/lmdb-go/v2/lmdbscan"
```

A utility package for scanning database ranges. The API is inspired by
[bufio.Scanner](https://godoc.org/bufio#Scanner) and the python cursor
[implementation](https://lmdb.readthedocs.org/en/release/#cursor-class).

#### exp/lmdbsync [![GoDoc](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/exp/lmdbsync?status.svg)](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/exp/lmdbsync) [![experimental](https://img.shields.io/badge/stability-experimental-red.svg)](#user-content-versioning-and-stability)


```go
import "github.com/PowerDNS/lmdb-go/v2/exp/lmdbsync"
```

An experimental utility package that provides synchronization necessary to
change an environment's map size after initialization.  The package provides
error handlers to automatically manage database size and retry failed
transactions.

The **lmdbsync** package is usable but the implementation of Handlers are
unstable and may change in incompatible ways without notice.  The use cases of
dynamic map sizes and multiprocessing are niche and the package requires much
more development driven by practical feedback before the Handler API and the
provided implementations can be considered stable.

## Key Features

### Idiomatic API

API inspired by [BoltDB](https://github.com/boltdb/bolt) with automatic
commit/rollback of transactions.  The goal of lmdb-go is to provide idiomatic
database interactions without compromising the flexibility of the C API.

**NOTE:** While the lmdb package tries hard to make LMDB as easy to use as
possible there are compromises, gotchas, and caveats that application
developers must be aware of when relying on LMDB to store their data.  All
users are encouraged to fully read the
[documentation](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdb) so they are
aware of these caveats.

Where the lmdb package and its implementation decisions do not meet the needs
of application developers in terms of safety or operational use the lmdbsync
package has been designed to wrap lmdb and safely fill in additional
functionality.  Consult the
[documentation](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/exp/lmdbsync) for
more information about the lmdbsync package.

### API coverage

The lmdb-go project aims for complete coverage of the LMDB C API (within
reason).  Some notable features and optimizations that are supported:

- Idiomatic subtransactions ("sub-updates") that allow the batching of updates.

- Batch IO on databases utilizing the `MDB_DUPSORT` and `MDB_DUPFIXED` flags.

- Reserved writes than can save in memory copies converting/buffering into
  `[]byte`.

For tracking purposes a list of unsupported features is kept in an
[issue](https://github.com/PowerDNS/lmdb-go/issues/1).

### Zero-copy reads

Applications with high performance requirements can opt-in to fast, zero-copy
reads at the cost of runtime safety.  Zero-copy behavior is specified at the
transaction level to reduce instrumentation overhead.

```
err := lmdb.View(func(txn *lmdb.Txn) error {
    // RawRead enables zero-copy behavior with some serious caveats.
    // Read the documentation carefully before using.
    txn.RawRead = true

    val, err := txn.Get(dbi, []byte("largevalue"), 0)
    // ...
})
```

### Documentation

Comprehensive documentation and examples are provided to demonstrate safe usage
of lmdb.  In addition to [godoc](https://godoc.org/github.com/PowerDNS/lmdb-go)
documentation, implementations of the standand LMDB commands (`mdb_stat`, etc)
can be found in the [cmd/](cmd/) directory and some simple experimental
commands can be found in the [exp/cmd/](exp/cmd) directory.  Aside from
providing minor utility these programs are provided as examples of lmdb in
practice.

## LMDB compared to BoltDB

BoltDB is a quality database with a design similar to LMDB.  Both store
key-value data in a file and provide ACID transactions.  So there are often
questions of why to use one database or the other.

### Advantages of BoltDB

- Nested databases allow for hierarchical data organization.

- Far more databases can be accessed concurrently.

- Operating systems that do not support sparse files do not use up excessive
  space due to a large pre-allocation of file space.  The exp/lmdbsync package
  is intended to resolve this problem with LMDB but it is not ready.

- As a pure Go package bolt can be easily cross-compiled using the `go`
  toolchain and `GOOS`/`GOARCH` variables.

- Its simpler design and implementation in pure Go mean it is free of many
  caveats and gotchas which are present using the lmdb package.  For more
  information about caveats with the lmdb package, consult its
  [documentation](https://godoc.org/github.com/PowerDNS/lmdb-go/v2/lmdb).

### Advantages of LMDB

- Keys can contain multiple values using the DupSort flag.

- Updates can have sub-updates for atomic batching of changes.

- Databases typically remain open for the application lifetime.  This limits
  the number of concurrently accessible databases.  But, this minimizes the
  overhead of database accesses and typically produces cleaner code than
  an equivalent BoltDB implementation.

- Significantly faster than BoltDB.  The raw speed of LMDB easily surpasses
  BoltDB.  Additionally, LMDB provides optimizations ranging from safe,
  feature-specific optimizations to generally unsafe, extremely situational
  ones.  Applications are free to enable any optimizations that fit their data,
  access, and reliability models.

- LMDB allows multiple applications to access a database simultaneously.
  Updates from concurrent processes are synchronized using a database lock
  file.

- As a C library, applications in any language can interact with LMDB
  databases.  Mission critical Go applications can use a database while Python
  scripts perform analysis on the side.

## Build

There is no dependency on shared libraries.  So most users can simply install
using `go get`.

`go get github.com/PowerDNS/lmdb-go/v2/lmdb`

On FreeBSD 10, you must explicitly set `CC` (otherwise it will fail with a
cryptic error), for example:

    CC=clang go test -v ./...

Building commands and running tests can be done with `go` or with `make`

    make bin
    make test
    make check
    make all

On Linux, you can specify the `pwritev` build tag to reduce the number of syscalls
required when committing a transaction. In your own package you can then do

    go build -tags pwritev .

to enable the optimisation.

## Documentation

### Go doc

The `go doc` documentation available on
[godoc.org](https://godoc.org/github.com/PowerDNS/lmdb-go) is the primary source
of developer documentation for lmdb-go.  It provides an overview of the API
with a lot of usage examples.  Where necessary the documentation points out
differences between the semantics of methods and their C counterparts.

### LMDB

The LMDB [homepage](http://symas.com/mdb/) and mailing list
([archives](http://www.openldap.org/lists/openldap-technical/)) are the
official source of documentation regarding low-level LMDB operation and
internals.

Along with an API reference LMDB provides a high-level
[summary](http://symas.com/mdb/doc/starting.html) of the library.  While
lmdb-go abstracts many of the thread and transaction details by default the
rest of the guide is still useful to compare with `go doc`.

### Versioning and Stability

The lmdb-go project makes regular releases with IDs `X.Y.Z`.  All packages
outside of the `exp/` directory are considered stable and adhere to the
guidelines of [semantic versioning](http://semver.org/).

Experimental packages (those packages in `exp/`) are not required to adhere to
semantic versioning.  However packages specifically declared to merely be
"unstable" can be relied on more for long term use with less concern.

The API of an unstable package may change in subtle ways between minor release
versions.  But deprecations will be indicated at least one release in advance
and all functionality will remain available through some method.

## License

Except where otherwise noted files in the lmdb-go project are licensed under
the BSD 3-clause open source license.

The LMDB C source is licensed under the OpenLDAP Public License.

## Links

#### [github.com/bmatsuo/raft-mdb](https://github.com/bmatsuo/raft-mdb) ([godoc](https://godoc.org/github.com/bmatsuo/raft-mdb))

An experimental backend for
[github.com/hashicorp/raft](https://github.com/hashicorp/raft) forked from
[github.com/hashicorp/raft-mdb](https://github.com/hashicorp/raft-mdb).

#### [github.com/bmatsuo/cayley/graph/lmdb](https://github.com/bmatsuo/cayley/tree/master/graph/lmdb) ([godoc](https://godoc.org/github.com/bmatsuo/cayley/graph/lmdb))

Experimental backend quad-store for
[github.com/google/cayley](https://github.com/google/cayley) based off of the
BoltDB
[implementation](https://github.com/google/cayley/tree/master/graph/bolt).
