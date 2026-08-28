package lmdb

import (
	"bytes"
	"errors"
	"testing"
)

func testOpenVersioned(t *testing.T, dir string, req Format) (*Env, error) {
	t.Helper()
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	if req.valid() {
		if err := env.SetFormat(req); err != nil {
			env.Close()
			t.Fatal(err)
		}
	}
	if err := env.Open(dir, 0, 0644); err != nil {
		env.Close()
		return nil, err
	}
	return env, nil
}

func TestEnv_Format(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()

	if got := env.Format(); got != FormatUnknown {
		t.Errorf("unopened env version: %v (!= FormatUnknown)", got)
	}
	if _, err := env.EngineVersion(); !errors.Is(err, ErrVersionUndetermined) {
		t.Errorf("expected ErrVersionUndetermined, got %v", err)
	}

	if err := env.Open(t.TempDir(), 0, 0644); err != nil {
		t.Fatal(err)
	}
	got := env.Format()
	if !got.valid() {
		t.Fatalf("opened env has no valid version: %v", got)
	}
	ev, err := env.EngineVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("engine under test: %s (Format %s)", ev, got)
	wantMajor := map[Format]int{V09: 0, V10: 1}[got]
	if ev.Major != wantMajor {
		t.Errorf("EngineVersion.Major = %d, want %d for %v", ev.Major, wantMajor, got)
	}
	if ev.Release == "" {
		t.Error("empty EngineVersion.Release")
	}
}

func TestEnv_SetFormat_afterOpen(t *testing.T) {
	env := setup(t)
	defer clean(env, t)
	if err := env.SetFormat(V09); err == nil {
		t.Error("expected error from SetFormat after Open")
	}
}

func TestEnv_SetFormat_invalid(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if err := env.SetFormat(Format(7)); !errors.Is(err, ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat, got %v", err)
	}
	if err := SetDefaultFormat(Format(7)); !errors.Is(err, ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat, got %v", err)
	}
}

// TestEnv_versionConflict: an explicitly requested version that contradicts
// an existing database's format must fail with ErrFormatConflict, in both
// directions.
func TestEnv_versionConflict(t *testing.T) {
	if !v10Available {
		t.Skip("LMDB 1.0 engine not available in this build")
	}
	for _, tc := range []struct{ have, want Format }{
		{V09, V10},
		{V10, V09},
	} {
		t.Run(tc.have.String()+"-then-"+tc.want.String(), func(t *testing.T) {
			dir := t.TempDir()
			createEnvFile(t, dir, tc.have, 0)
			_, err := testOpenVersioned(t, dir, tc.want)
			if !errors.Is(err, ErrFormatConflict) {
				t.Fatalf("expected ErrFormatConflict, got %v", err)
			}
		})
	}
}

// TestEnv_existingFormatWins: without an explicit request, an existing
// database's format beats the process default.
func TestEnv_existingFormatWins(t *testing.T) {
	if !v10Available {
		t.Skip("LMDB 1.0 engine not available in this build")
	}
	dir := t.TempDir()
	createEnvFile(t, dir, V09, 0)

	if err := SetDefaultFormat(V10); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := SetDefaultFormat(FormatUnknown); err != nil {
			t.Fatal(err)
		}
	}()

	env, err := testOpenVersioned(t, dir, FormatUnknown)
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close()
	if got := env.Format(); got != V09 {
		t.Errorf("existing 0.9 database opened as %v", got)
	}
}

func TestSetDefaultFormat(t *testing.T) {
	if !v10Available {
		t.Skip("LMDB 1.0 engine not available in this build")
	}
	// Explicit setter beats the environment variable.
	t.Setenv(DefaultFormatEnvVar, "09")
	if err := SetDefaultFormat(V10); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := SetDefaultFormat(FormatUnknown); err != nil {
			t.Fatal(err)
		}
	}()
	if got := defaultFormat(); got != V10 {
		t.Errorf("default version %v (!= V10, explicit setter must win)", got)
	}

	// Clearing the setter falls back to the environment variable.
	if err := SetDefaultFormat(FormatUnknown); err != nil {
		t.Fatal(err)
	}
	if got := defaultFormat(); got != V09 {
		t.Errorf("default version %v (!= V09 from env var)", got)
	}
	t.Setenv(DefaultFormatEnvVar, "10")
	if got := defaultFormat(); got != V10 {
		t.Errorf("default version %v (!= V10 from env var)", got)
	}
	t.Setenv(DefaultFormatEnvVar, "")
	if got := defaultFormat(); got != V09 {
		t.Errorf("default version %v (!= built-in V09)", got)
	}
}

// TestEnv_crossFormatSimultaneous opens a 0.9-format and a 1.0-format
// environment in the same process at the same time and reads/writes both -
// the core capability that motivated bundling both engines.
func TestEnv_crossFormatSimultaneous(t *testing.T) {
	if !v10Available {
		t.Skip("LMDB 1.0 engine not available in this build")
	}

	envs := map[Format]*Env{}
	for _, ver := range []Format{V09, V10} {
		env, err := testOpenVersioned(t, t.TempDir(), ver)
		if err != nil {
			t.Fatal(err)
		}
		defer env.Close()
		if got := env.Format(); got != ver {
			t.Fatalf("requested %v, got %v", ver, got)
		}
		envs[ver] = env
	}

	// Write a distinct value in each env, then read both back.
	for ver, env := range envs {
		val := []byte("value-" + ver.String())
		err := env.Update(func(txn *Txn) error {
			db, err := txn.OpenRoot(0)
			if err != nil {
				return err
			}
			return txn.Put(db, []byte("key"), val, 0)
		})
		if err != nil {
			t.Fatalf("%v: %v", ver, err)
		}
	}
	for ver, env := range envs {
		want := []byte("value-" + ver.String())
		err := env.View(func(txn *Txn) error {
			db, err := txn.OpenRoot(0)
			if err != nil {
				return err
			}
			got, err := txn.Get(db, []byte("key"))
			if err != nil {
				return err
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%v: got %q (!= %q)", ver, got, want)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("%v: %v", ver, err)
		}

		ev, err := env.EngineVersion()
		if err != nil {
			t.Fatalf("%v: %v", ver, err)
		}
		wantMajor := map[Format]int{V09: 0, V10: 1}[ver]
		if ev.Major != wantMajor {
			t.Errorf("%v: engine major %d (!= %d)", ver, ev.Major, wantMajor)
		}
	}

	// Both databases sniff back to their formats after close-independent
	// reads (paths still open here; sniffing is read-only and safe).
	for ver, env := range envs {
		path, err := env.Path()
		if err != nil {
			t.Fatal(err)
		}
		sniffed, err := SniffFormat(path, 0)
		if err != nil {
			t.Fatal(err)
		}
		if sniffed != ver {
			t.Errorf("sniffed %v (!= %v)", sniffed, ver)
		}
	}
}

// TestEnv_PrevSnapshot exercises the lmdb.PrevSnapshot open flag: on the 1.0
// engine it opens the previous snapshot (losing the latest transaction); the
// 0.9 engine has no equivalent and rejects the flag with EINVAL.
func TestEnv_PrevSnapshot(t *testing.T) {
	key := []byte("psnap-key")
	put := func(env *Env, val string) error {
		return env.Update(func(txn *Txn) error {
			db, err := txn.OpenRoot(0)
			if err != nil {
				return err
			}
			return txn.Put(db, key, []byte(val), 0)
		})
	}

	t.Run("v10", func(t *testing.T) {
		if !v10Available {
			t.Skip("LMDB 1.0 engine not available in this build")
		}
		dir := t.TempDir()
		env, err := testOpenVersioned(t, dir, V10)
		if err != nil {
			t.Fatal(err)
		}
		if err := put(env, "first"); err != nil {
			t.Fatal(err)
		}
		if err := put(env, "second"); err != nil {
			t.Fatal(err)
		}
		if err := env.Close(); err != nil {
			t.Fatal(err)
		}

		prev, err := NewEnv()
		if err != nil {
			t.Fatal(err)
		}
		defer prev.Close()
		if err := prev.Open(dir, PrevSnapshot|Readonly, 0644); err != nil {
			t.Fatal(err)
		}
		err = prev.View(func(txn *Txn) error {
			db, err := txn.OpenRoot(0)
			if err != nil {
				return err
			}
			got, err := txn.Get(db, key)
			if err != nil {
				return err
			}
			if string(got) != "first" {
				t.Errorf("previous snapshot has %q (want %q)", got, "first")
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("v09rejects", func(t *testing.T) {
		dir := t.TempDir()
		env, err := testOpenVersioned(t, dir, V09)
		if err != nil {
			t.Fatal(err)
		}
		if err := env.Close(); err != nil {
			t.Fatal(err)
		}

		env2, err := NewEnv()
		if err != nil {
			t.Fatal(err)
		}
		defer env2.Close()
		err = env2.Open(dir, PrevSnapshot|Readonly, 0644)
		if !errIsEINVAL(err) {
			t.Errorf("expected EINVAL from the 0.9 engine, got %v", err)
		}
	})
}
