// Package coexist verifies that lmdb-go v1 and v2 can be imported by the
// same binary. This is a separate Go module (see go.mod here) so the main
// module never depends on v1.
//
// Linking this test binary at all proves there are no duplicate C symbols:
// v1 carries an unprefixed LMDB 0.9 plus lmdbgo_* helpers, while all v2 C
// symbols are namespaced (mdb09_*/mdb10_* engines, lmdbgo2_* glue). The
// tests then prove the two modules interoperate on the same database files.
package coexist

import (
	"bytes"
	"testing"

	lmdb1 "github.com/PowerDNS/lmdb-go/lmdb"
	lmdb2 "github.com/PowerDNS/lmdb-go/v2/lmdb"
)

var (
	key = []byte("coexist-key")
	val = []byte("coexist-val")
)

// TestV1WriteV2Read writes a database with v1 (LMDB 0.9 format) and reads it
// back with v2, which must sniff it as the 0.9 format.
func TestV1WriteV2Read(t *testing.T) {
	dir := t.TempDir()

	env1, err := lmdb1.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	if err := env1.Open(dir, 0, 0644); err != nil {
		t.Fatal(err)
	}
	err = env1.Update(func(txn *lmdb1.Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		return txn.Put(db, key, val, 0)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := env1.Close(); err != nil {
		t.Fatal(err)
	}

	if ver, err := lmdb2.SniffFormat(dir, 0); err != nil || ver != lmdb2.V09 {
		t.Fatalf("v2 sniffed v1 database as %v, %v", ver, err)
	}

	env2, err := lmdb2.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env2.Close()
	if err := env2.Open(dir, 0, 0644); err != nil {
		t.Fatal(err)
	}
	if got := env2.Format(); got != lmdb2.V09 {
		t.Errorf("v2 opened v1 database with engine %v", got)
	}
	err = env2.View(func(txn *lmdb2.Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		got, err := txn.Get(db, key)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, val) {
			t.Errorf("got %q (!= %q)", got, val)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestV2WriteV1Read writes with v2 (default 0.9 format) and reads with v1.
func TestV2WriteV1Read(t *testing.T) {
	dir := t.TempDir()

	env2, err := lmdb2.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	if err := env2.Open(dir, 0, 0644); err != nil {
		t.Fatal(err)
	}
	if got := env2.Format(); got != lmdb2.V09 {
		t.Fatalf("v2 default format is %v (!= V09)", got)
	}
	err = env2.Update(func(txn *lmdb2.Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		return txn.Put(db, key, val, 0)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := env2.Close(); err != nil {
		t.Fatal(err)
	}

	env1, err := lmdb1.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env1.Close()
	if err := env1.Open(dir, 0, 0644); err != nil {
		t.Fatal(err)
	}
	err = env1.View(func(txn *lmdb1.Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		got, err := txn.Get(db, key)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, val) {
			t.Errorf("got %q (!= %q)", got, val)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestV1CannotOpenV10Format documents why v2 exists: a 1.0-format database
// created by v2 is not readable by v1 (its LMDB 0.9 fails the open), while
// v2 sniffs and opens it fine.
func TestV1CannotOpenV10Format(t *testing.T) {
	dir := t.TempDir()

	env2, err := lmdb2.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	if err := env2.SetFormat(lmdb2.V10); err != nil {
		t.Fatal(err)
	}
	if err := env2.Open(dir, 0, 0644); err != nil {
		t.Fatal(err)
	}
	err = env2.Update(func(txn *lmdb2.Txn) error {
		db, err := txn.OpenRoot(0)
		if err != nil {
			return err
		}
		return txn.Put(db, key, val, 0)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := env2.Close(); err != nil {
		t.Fatal(err)
	}

	if ver, err := lmdb2.SniffFormat(dir, 0); err != nil || ver != lmdb2.V10 {
		t.Fatalf("v2 sniffed as %v, %v", ver, err)
	}

	env1, err := lmdb1.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env1.Close()
	if err := env1.Open(dir, 0, 0644); err == nil {
		t.Error("v1 unexpectedly opened a 1.0-format database")
	} else {
		t.Logf("v1 open of a 1.0-format database fails as expected: %v", err)
	}
}
