package migrate

import (
	"os"
	"path/filepath"
	"testing"
)

// chdirWithTemplates runs the test from an empty directory that has the
// repository's template/ directory in it, which is what genFile reads from
// the working directory.
func chdirWithTemplates(t *testing.T) string {
	t.Helper()
	tmpl, err := filepath.Abs(filepath.Join("..", "..", "template"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Symlink(tmpl, filepath.Join(dir, "template")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	return dir
}

func TestGenFileWritesTheMigration(t *testing.T) {
	dir := chdirWithTemplates(t)
	out := filepath.Join(dir, "cmd", "migrate", "migration", "version-local")
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := genFile(); err != nil {
		t.Fatalf("genFile: %v", err)
	}
	files, err := filepath.Glob(filepath.Join(out, "*_migrate.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("generated files = %v, want one", files)
	}
}

// Run from a directory without cmd/migrate/migration, writing the file fails.
// The failure used to end the process from inside pkg.FileCreate and, before
// that, was dropped by genFile's caller; it now reaches the caller.
func TestGenFileReportsAFileItCannotWrite(t *testing.T) {
	chdirWithTemplates(t)

	if err := genFile(); err == nil {
		t.Fatal("genFile returned nil with no directory to write into")
	}
}
