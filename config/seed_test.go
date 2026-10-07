package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSeedFallsBackToEmbedded(t *testing.T) {
	// A directory with no config/ in it, as in a container.
	missing := filepath.Join(t.TempDir(), "config", "db.sql")
	b, err := ReadSeed(missing)
	if err != nil || len(b) == 0 {
		t.Fatalf("embedded fallback: len=%d err=%v", len(b), err)
	}
}

func TestReadSeedPrefersDisk(t *testing.T) {
	p := filepath.Join(t.TempDir(), "db.sql")
	if err := os.WriteFile(p, []byte("-- mine;"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := ReadSeed(p)
	if err != nil || string(b) != "-- mine;" {
		t.Fatalf("got %q, %v", b, err)
	}
}

func TestReadSeedUnknownFileStillErrors(t *testing.T) {
	if _, err := ReadSeed(filepath.Join(t.TempDir(), "nope.sql")); err == nil {
		t.Fatal("expected an error for a script that is not embedded")
	}
}
