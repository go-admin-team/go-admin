package main

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

// The whole 008 chain, against a binary that was actually built and a
// database that was actually migrated.
//
// Everything below this level is covered by unit tests with an injected
// engine and a hand-built schema, which is where the shapes are pinned down.
// What only this can catch is the wiring: that an application's init()
// reaches both registries, that `install` finds a manifest through
// app.Snapshot, that the seeder writes what the uninstaller looks for, and
// that the command exits non-zero when a migration fails - the last of which
// a deployment reads to decide whether to start the new version.
const settings = `settings:
  application:
    host: 0.0.0.0
    mode: dev
    name: e2e
    port: 8000
    readtimeout: 10000
    writertimeout: 20000
  database:
    driver: sqlite3
    source: ./e2e.db
  jwt:
    secret: e2e
    timeout: 3600
  logger:
    path: temp/logs
    stdout: default
    level: error
    enableddb: false
  queue:
    memory:
      poolSize: 10
`

type env struct {
	t   *testing.T
	dir string
	bin string
}

// newEnv builds the binary and lays out a working directory for it.
func newEnv(t *testing.T) *env {
	t.Helper()
	if testing.Short() {
		t.Skip("builds a binary and migrates a database; skipped under -short")
	}

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "temp", "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config", "settings.yml"), []byte(settings), 0o644); err != nil {
		t.Fatal(err)
	}
	// The framework's first migration reads this file rather than carrying
	// the rows in Go.
	seedSQL, err := os.ReadFile(filepath.Join("..", "..", "config", "db.sql"))
	if err != nil {
		t.Fatalf("reading config/db.sql: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config", "db.sql"), seedSQL, 0o644); err != nil {
		t.Fatal(err)
	}

	bin := filepath.Join(dir, "go-admin-e2e")
	build := exec.Command("go", "build", "-tags", "sqlite3", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building the binary: %v\n%s", err, out)
	}
	return &env{t: t, dir: dir, bin: bin}
}

// run executes the binary and returns its combined output and exit code.
func (e *env) run(args ...string) (string, int) {
	e.t.Helper()
	args = append(args, "-c", "config/settings.yml")
	cmd := exec.Command(e.bin, args...)
	cmd.Dir = e.dir
	out, err := cmd.CombinedOutput()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		e.t.Fatalf("running %v: %v", args, err)
	}
	return string(out), code
}

func (e *env) mustRun(args ...string) string {
	e.t.Helper()
	out, code := e.run(args...)
	if code != 0 {
		e.t.Fatalf("%v exited %d:\n%s", args, code, out)
	}
	return out
}

// exec runs one statement against the database the binary uses.
func (e *env) exec(stmt string) {
	e.t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(e.dir, "e2e.db"))
	if err != nil {
		e.t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(stmt); err != nil {
		e.t.Fatalf("%s: %v", stmt, err)
	}
}

func (e *env) count(query string, args ...any) int {
	e.t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(e.dir, "e2e.db"))
	if err != nil {
		e.t.Fatal(err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		e.t.Fatalf("%s: %v", query, err)
	}
	return n
}

func TestInstallUninstallReinstall(t *testing.T) {
	e := newEnv(t)

	// Framework only. The application's migration must not run here, or the
	// install below has nothing left to do and the interesting half of it
	// goes untested - which is what happens if this uses plain `migrate`.
	e.mustRun("migrate", "--app", "core")
	if n := e.count("SELECT COUNT(*) FROM sys_migration WHERE app_code = 'order'"); n != 0 {
		t.Fatalf("the application's migration ran during the framework's: %d rows", n)
	}
	if n := e.count("SELECT COUNT(*) FROM sqlite_master WHERE name = 'app_order'"); n != 0 {
		t.Fatal("the application's own table exists before it was installed")
	}

	// A1.
	out := e.mustRun("migrate", "install", "order")
	if !strings.Contains(out, "order-1793800000000") {
		t.Errorf("the install did not report applying the migration:\n%s", out)
	}
	// A7.
	if !strings.Contains(out, "rebuild") {
		t.Errorf("the install did not say the code is not running yet:\n%s", out)
	}
	for _, c := range []struct {
		what  string
		query string
		want  int
	}{
		{"menus", "SELECT COUNT(*) FROM sys_menu WHERE app_code = 'order'", 4},
		{"apis", "SELECT COUNT(*) FROM sys_api WHERE app_code = 'order'", 4},
		{"ledger", "SELECT COUNT(*) FROM sys_app_casbin_grant WHERE app_code = 'order'", 4},
		{"policies", "SELECT COUNT(*) FROM casbin_rule WHERE v1 LIKE '/api/v1/order%'", 4},
		{"migration records", "SELECT COUNT(*) FROM sys_migration WHERE app_code = 'order'", 1},
		{"sys_app rows", "SELECT COUNT(*) FROM sys_app WHERE app_code = 'order'", 1},
		{"its own table", "SELECT COUNT(*) FROM sqlite_master WHERE name = 'app_order'", 1},
		{"status=installed", "SELECT COUNT(*) FROM sys_app WHERE app_code = 'order' AND status = 2", 1},
	} {
		if n := e.count(c.query); n != c.want {
			t.Errorf("after install, %s = %d, want %d", c.what, n, c.want)
		}
	}

	// A2: installing the same version again does nothing and says so.
	out = e.mustRun("migrate", "install", "order")
	if !strings.Contains(out, "already installed") {
		t.Errorf("a second install was not reported as a no-op:\n%s", out)
	}

	// A3: the application's own data is not the uninstaller's to remove.
	e.exec("INSERT INTO app_order (created_at, updated_at) VALUES (datetime('now'), datetime('now'))")
	if n := e.count("SELECT COUNT(*) FROM app_order"); n != 1 {
		t.Fatalf("the business row was not written: %d", n)
	}

	out = e.mustRun("migrate", "uninstall", "order")
	if !strings.Contains(out, "own tables were not touched") {
		t.Errorf("the uninstall did not say what it left alone:\n%s", out)
	}
	for _, c := range []struct {
		what  string
		query string
	}{
		{"menus", "SELECT COUNT(*) FROM sys_menu WHERE app_code = 'order'"},
		{"apis", "SELECT COUNT(*) FROM sys_api WHERE app_code = 'order'"},
		{"ledger", "SELECT COUNT(*) FROM sys_app_casbin_grant WHERE app_code = 'order'"},
		{"policies", "SELECT COUNT(*) FROM casbin_rule WHERE v1 LIKE '/api/v1/order%'"},
		{"migration records", "SELECT COUNT(*) FROM sys_migration WHERE app_code = 'order'"},
		{"sys_app rows", "SELECT COUNT(*) FROM sys_app WHERE app_code = 'order'"},
	} {
		if n := e.count(c.query); n != 0 {
			t.Errorf("after uninstall, %s = %d, want 0", c.what, n)
		}
	}
	if n := e.count("SELECT COUNT(*) FROM sqlite_master WHERE name = 'app_order'"); n != 1 {
		t.Error("the uninstall dropped the application's own table")
	}
	if n := e.count("SELECT COUNT(*) FROM app_order"); n != 1 {
		t.Errorf("the uninstall removed %d business row(s)", 1-n)
	}

	// A4: the migration records had to go, or this reinstall finds every
	// version applied, seeds nothing, and reports success.
	e.mustRun("migrate", "install", "order")
	for _, c := range []struct {
		what  string
		query string
		want  int
	}{
		{"menus", "SELECT COUNT(*) FROM sys_menu WHERE app_code = 'order'", 4},
		{"ledger", "SELECT COUNT(*) FROM sys_app_casbin_grant WHERE app_code = 'order'", 4},
		{"policies", "SELECT COUNT(*) FROM casbin_rule WHERE v1 LIKE '/api/v1/order%'", 4},
	} {
		if n := e.count(c.query); n != c.want {
			t.Errorf("after reinstall, %s = %d, want %d", c.what, n, c.want)
		}
	}
	if n := e.count("SELECT COUNT(*) FROM app_order"); n != 1 {
		t.Error("the business row did not survive an uninstall and reinstall")
	}
}

// A deployment decides whether to start the new version on this exit code.
func TestMigrateExitsNonZeroWhenItFails(t *testing.T) {
	e := newEnv(t)

	// The framework's first migration reads config/db.sql. Without it the
	// migration fails, which is the cheapest real failure to arrange.
	if err := os.Remove(filepath.Join(e.dir, "config", "db.sql")); err != nil {
		t.Fatal(err)
	}
	out, code := e.run("migrate")
	if code == 0 {
		t.Errorf("a failed migration exited 0:\n%s", out)
	}
}

func TestUnknownAppIsRefused(t *testing.T) {
	e := newEnv(t)
	e.mustRun("migrate", "--app", "core")

	out, code := e.run("migrate", "install", "ordr")
	if code == 0 {
		t.Errorf("a mistyped code was installed:\n%s", out)
	}
	if !strings.Contains(out, "order") {
		t.Errorf("the refusal does not name what is registered:\n%s", out)
	}
}
