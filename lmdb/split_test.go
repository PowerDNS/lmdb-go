package lmdb

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

// freelistEntries returns the number of records in the freelist DB (DBI 0).
func freelistEntries(t *testing.T, env *Env) uint64 {
	t.Helper()
	var n uint64
	err := env.View(func(txn *Txn) error {
		stat, err := txn.Stat(DBI(0))
		if err != nil {
			return err
		}
		n = stat.Entries
		return nil
	})
	if err != nil {
		t.Fatalf("freelist stat: %v", err)
	}
	return n
}

// TestEnv_SetSplitFreelist frees more pages in one transaction than one
// freelist record can hold, and checks that the record is split by default
// and when enabled explicitly, but not when turned off, and that CommitID
// reports the resulting ID.
func TestEnv_SetSplitFreelist(t *testing.T) {
	cases := []struct {
		name  string
		set   bool // call SetSplitFreelist
		split bool
	}{
		{"default", false, true},
		{"off", true, false},
		{"on", true, true},
	}
	for _, c := range cases {
		split := c.split
		t.Run(c.name, func(t *testing.T) {
			env := setup(t)
			defer clean(env, t)
			if err := env.SetMapSize(512 << 20); err != nil {
				t.Fatal(err)
			}
			if c.set {
				if err := env.SetSplitFreelist(split); err != nil {
					t.Fatal(err)
				}
			}
			stat, err := env.Stat()
			if err != nil {
				t.Fatal(err)
			}

			// A value of one page takes two overflow pages.  One freelist
			// record holds (psize-16)/8-1 page numbers: 509 with 4 KB pages,
			// 2045 with 16 KB pages.
			const n = 1500
			val := bytes.Repeat([]byte("x"), int(stat.PSize))
			key := func(i int) []byte { return []byte(fmt.Sprintf("key%06d", i)) }
			var dbi DBI
			err = env.Update(func(txn *Txn) (err error) {
				dbi, err = txn.OpenRoot(0)
				if err != nil {
					return err
				}
				for i := 0; i < n; i++ {
					if err := txn.Put(dbi, key(i), val, 0); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				t.Fatalf("fill: %v", err)
			}
			before := freelistEntries(t, env)

			var del *Txn
			var idInside uintptr
			err = env.Update(func(txn *Txn) error {
				del = txn
				idInside = txn.ID()
				for i := 0; i < n; i++ {
					if err := txn.Del(dbi, key(i), nil); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				t.Fatalf("delete: %v", err)
			}
			after := freelistEntries(t, env)
			info, err := env.Info()
			if err != nil {
				t.Fatal(err)
			}
			commitID := del.CommitID()
			t.Logf("psize=%d freelist records before=%d after=%d id=%d commit=%d",
				stat.PSize, before, after, idInside, commitID)

			if commitID != uintptr(info.LastTxnID) {
				t.Errorf("CommitID %d != LastTxnID %d", commitID, info.LastTxnID)
			}
			if split {
				if after-before < 2 {
					t.Errorf("freelist record not split: %d record(s) added", after-before)
				}
				if commitID <= idInside {
					t.Errorf("CommitID %d not past ID %d after a split", commitID, idInside)
				}
			} else {
				if after-before != 1 {
					t.Errorf("expected one freelist record, %d added", after-before)
				}
				if commitID != idInside {
					t.Errorf("CommitID %d != ID %d", commitID, idInside)
				}
			}

			// The data must be gone and the environment still usable.
			err = env.Update(func(txn *Txn) error {
				if _, err := txn.Get(dbi, key(0)); !IsNotFound(err) {
					return fmt.Errorf("key0 after delete: %v", err)
				}
				return txn.Put(dbi, key(0), val, 0)
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestTxn_CommitID checks the CommitID contract for transactions that do not
// split anything: with changes, empty, nested, read-only and aborted.
func TestTxn_CommitID(t *testing.T) {
	env := setup(t)
	defer clean(env, t)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	txn, err := env.BeginTxn(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if txn.CommitID() != 0 {
		t.Errorf("CommitID before commit: %d", txn.CommitID())
	}
	dbi, err := txn.OpenRoot(0)
	if err != nil {
		t.Fatal(err)
	}
	if err := txn.Put(dbi, []byte("k"), []byte("v"), 0); err != nil {
		t.Fatal(err)
	}
	id := txn.ID()
	if err := txn.Commit(); err != nil {
		t.Fatal(err)
	}
	if txn.CommitID() != id {
		t.Errorf("write txn: CommitID %d != ID %d", txn.CommitID(), id)
	}

	// Nothing written: the previously committed ID.
	txn, err = env.BeginTxn(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if txn.ID() != id+1 {
		t.Fatalf("empty txn ID %d, expected %d", txn.ID(), id+1)
	}
	if err := txn.Commit(); err != nil {
		t.Fatal(err)
	}
	if txn.CommitID() != id {
		t.Errorf("empty txn: CommitID %d != previous ID %d", txn.CommitID(), id)
	}

	// Nested: the parent's ID, which is also the parent's outcome.
	parent, err := env.BeginTxn(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	child, err := env.BeginTxn(parent, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := child.Put(dbi, []byte("k2"), []byte("v"), 0); err != nil {
		t.Fatal(err)
	}
	if err := child.Commit(); err != nil {
		t.Fatal(err)
	}
	if child.CommitID() != parent.ID() {
		t.Errorf("child: CommitID %d != parent ID %d", child.CommitID(), parent.ID())
	}
	if err := parent.Commit(); err != nil {
		t.Fatal(err)
	}
	if parent.CommitID() != id+1 {
		t.Errorf("parent: CommitID %d != %d", parent.CommitID(), id+1)
	}

	// Aborted: nothing stored.
	txn, err = env.BeginTxn(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := txn.Put(dbi, []byte("k3"), []byte("v"), 0); err != nil {
		t.Fatal(err)
	}
	txn.Abort()
	if txn.CommitID() != 0 {
		t.Errorf("aborted txn: CommitID %d", txn.CommitID())
	}

	// Read-only: always 0.
	txn, err = env.BeginTxn(nil, Readonly)
	if err != nil {
		t.Fatal(err)
	}
	if err := txn.Commit(); err != nil {
		t.Fatal(err)
	}
	if txn.CommitID() != 0 {
		t.Errorf("read-only txn: CommitID %d", txn.CommitID())
	}
}
