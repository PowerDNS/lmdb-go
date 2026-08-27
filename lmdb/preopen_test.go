package lmdb

// Tests for the v2 pre-open behavior: NewEnv no longer creates the engine
// environment (the engine is only known at Open), so setters buffer on the
// Go side and getters emulate the engines' documented pre-open behavior.
// Where v1 crashed on an unopened env, v2 returns clean errors.

import (
	"errors"
	"syscall"
	"testing"
)

func TestEnv_preOpen_flags(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()

	if flags, err := env.Flags(); err != nil || flags != 0 {
		t.Errorf("fresh env flags: %d, %v", flags, err)
	}
	if err := env.SetFlags(NoSync | NoMetaSync); err != nil {
		t.Fatal(err)
	}
	if err := env.UnsetFlags(NoMetaSync); err != nil {
		t.Fatal(err)
	}
	flags, err := env.Flags()
	if err != nil {
		t.Fatal(err)
	}
	if flags != NoSync {
		t.Errorf("buffered flags %#x (!= NoSync)", flags)
	}
	// Non-changeable flags are rejected with EINVAL like the engines do.
	if err := env.SetFlags(WriteMap); !errIsEINVAL(err) {
		t.Errorf("expected EINVAL for WriteMap pre-open, got %v", err)
	}

	// The buffered flag is replayed onto the engine env at Open.
	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	flags, err = env.Flags()
	if err != nil {
		t.Fatal(err)
	}
	if flags&NoSync == 0 {
		t.Errorf("NoSync not replayed at Open: flags %#x", flags)
	}
}

func errIsEINVAL(err error) bool {
	var op *OpError
	return errors.As(err, &op) && op.Errno == syscall.EINVAL
}

func TestEnv_preOpen_maxReaders(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()

	if n, err := env.MaxReaders(); err != nil || n != 126 {
		t.Errorf("default MaxReaders: %d, %v (want 126)", n, err)
	}
	if err := env.SetMaxReaders(150); err != nil {
		t.Fatal(err)
	}
	if n, err := env.MaxReaders(); err != nil || n != 150 {
		t.Errorf("buffered MaxReaders: %d, %v (want 150)", n, err)
	}
	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	if n, err := env.MaxReaders(); err != nil || n != 150 {
		t.Errorf("replayed MaxReaders: %d, %v (want 150)", n, err)
	}
	// Post-open the engines reject the setter with EINVAL, as in v1.
	if err := env.SetMaxReaders(200); !errIsEINVAL(err) {
		t.Errorf("expected EINVAL post-open, got %v", err)
	}
}

func TestEnv_preOpen_maxKeySize(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if n := env.MaxKeySize(); n != 511 {
		t.Errorf("pre-open MaxKeySize %d (!= 511)", n)
	}
	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	// Post-open the value is engine-dependent: 0.9's compile-time constant
	// is 511; 1.0 computes it from the page size (8122 for 4K pages).
	n := env.MaxKeySize()
	switch env.LMDBVersion() {
	case V09:
		if n != 511 {
			t.Errorf("post-open MaxKeySize %d (!= 511 on the 0.9 engine)", n)
		}
	case V10:
		if n < 511 {
			t.Errorf("post-open MaxKeySize %d (< 511 on the 1.0 engine)", n)
		}
		t.Logf("1.0 engine MaxKeySize: %d", n)
	}
}

func TestEnv_preOpen_errors(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()

	if _, err := env.Stat(); !errors.Is(err, errNotOpen) {
		t.Errorf("Stat: %v", err)
	}
	if _, err := env.Info(); !errors.Is(err, errNotOpen) {
		t.Errorf("Info: %v", err)
	}
	if err := env.Sync(true); !errors.Is(err, errNotOpen) {
		t.Errorf("Sync: %v", err)
	}
	if err := env.Copy(t.TempDir()); !errors.Is(err, errNotOpen) {
		t.Errorf("Copy: %v", err)
	}
	if _, err := env.Path(); !errors.Is(err, errNotOpen) {
		t.Errorf("Path: %v", err)
	}
	if _, err := env.FD(); !errors.Is(err, errNotOpen) {
		t.Errorf("FD: %v", err)
	}
	if _, err := env.BeginTxn(nil, 0); !errors.Is(err, errNotOpen) {
		t.Errorf("BeginTxn: %v", err)
	}
	err = env.View(func(txn *Txn) error { return nil })
	if !errors.Is(err, errNotOpen) {
		t.Errorf("View: %v", err)
	}
	env.CloseDBI(1) // must not crash

	if err := env.ReaderList(func(msg string) error {
		if msg != "(no reader locks)\n" {
			t.Errorf("unexpected pre-open reader list message %q", msg)
		}
		return nil
	}); err != nil {
		t.Errorf("ReaderList: %v", err)
	}
	if n, err := env.ReaderCheck(); err != nil || n != 0 {
		t.Errorf("ReaderCheck: %d, %v", n, err)
	}
}

func TestEnv_preOpen_replay(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if err := env.SetMapSize(16 << 20); err != nil {
		t.Fatal(err)
	}
	if err := env.SetMaxDBs(4); err != nil {
		t.Fatal(err)
	}
	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	info, err := env.Info()
	if err != nil {
		t.Fatal(err)
	}
	if info.MapSize != 16<<20 {
		t.Errorf("map size %d not replayed", info.MapSize)
	}
	// The buffered SetMaxDBs(4) must allow opening named databases.
	err = env.Update(func(txn *Txn) error {
		_, err := txn.OpenDBI("named", Create)
		return err
	})
	if err != nil {
		t.Error(err)
	}
	// SetMapSize keeps working post-open (the lmdbsync resize path).
	if err := env.SetMapSize(32 << 20); err != nil {
		t.Errorf("post-open SetMapSize: %v", err)
	}
}

func TestEnv_openTwice(t *testing.T) {
	env := setup(t)
	defer clean(env, t)
	if err := env.Open(t.TempDir(), 0, 0644); !errIsEINVAL(err) {
		t.Errorf("expected EINVAL for double Open, got %v", err)
	}
}

func TestEnv_openRetry(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if err := env.SetMaxDBs(2); err != nil {
		t.Fatal(err)
	}
	// Opening a nonexistent directory fails inside the engine...
	err = env.Open("/nonexistent-lmdb-go-test-dir/db", 0, 0644)
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
	// ...but the env stays usable: Open can be retried and the buffered
	// settings are replayed on the retry.
	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	err = env.Update(func(txn *Txn) error {
		_, err := txn.OpenDBI("named", Create)
		return err
	})
	if err != nil {
		t.Error(err)
	}
}

func TestEnv_closeBeforeOpen(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	if err := env.Close(); err != nil {
		t.Errorf("Close before Open: %v", err)
	}
	if err := env.Close(); err == nil {
		t.Error("expected error from double Close")
	}
	if err := env.Open(t.TempDir(), 0, 0644); !errors.Is(err, errClosed) {
		t.Errorf("Open after Close: %v", err)
	}
	if err := env.SetMapSize(1 << 20); !errors.Is(err, errClosed) {
		t.Errorf("SetMapSize after Close: %v", err)
	}
}
