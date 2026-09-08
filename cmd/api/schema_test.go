package api

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	commonmodels "go-admin/common/models"
)

func TestPendingVersionsReportsOnlyWhatIsNotApplied(t *testing.T) {
	registered := []string{"1000_a", "2000_b", "3000_c"}
	applied := map[string]bool{"1000_a": true, "3000_c": true}

	got := pendingVersions(registered, applied)
	if len(got) != 1 || got[0] != "2000_b" {
		t.Errorf("pending = %v, want [2000_b]", got)
	}
}

func TestPendingVersionsIsEmptyWhenTheDatabaseIsCurrent(t *testing.T) {
	registered := []string{"1000_a", "2000_b"}
	applied := map[string]bool{"1000_a": true, "2000_b": true}

	if got := pendingVersions(registered, applied); len(got) != 0 {
		t.Errorf("pending = %v, want none", got)
	}
}

// A row recorded that this binary no longer registers is not pending. It is
// the orphan `migrate status` already reports, and readiness has nothing to
// say about it: the schema is ahead, not behind, and requests will be served
// correctly.
func TestPendingVersionsIgnoresAppliedRowsNothingRegisters(t *testing.T) {
	registered := []string{"1000_a"}
	applied := map[string]bool{"1000_a": true, "9999_gone": true}

	if got := pendingVersions(registered, applied); len(got) != 0 {
		t.Errorf("pending = %v, want none - an orphaned row is not a pending migration", got)
	}
}

func memoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Migrator().DropTable(&commonmodels.Migration{}) })
	return db
}

// A first deploy has no sys_migration table. That is "nothing applied", not an
// error: reporting it as one would make the check fail for a reason the
// operator cannot act on, on the one deployment where every migration really
// is pending.
func TestAppliedVersionsTreatsAMissingTableAsNothingApplied(t *testing.T) {
	db := memoryDB(t)
	db.Migrator().DropTable(&commonmodels.Migration{})

	got, err := appliedVersions(context.Background(), db)
	if err != nil {
		t.Fatalf("appliedVersions: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("applied = %v, want empty", got)
	}
}

func TestAppliedVersionsReadsWhatTheTableHolds(t *testing.T) {
	db := memoryDB(t)
	if err := db.AutoMigrate(&commonmodels.Migration{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	db.Create(&commonmodels.Migration{Version: "1000_a"})
	db.Create(&commonmodels.Migration{Version: "2000_b"})

	got, err := appliedVersions(context.Background(), db)
	if err != nil {
		t.Fatalf("appliedVersions: %v", err)
	}
	if !got["1000_a"] || !got["2000_b"] || len(got) != 2 {
		t.Errorf("applied = %v, want the two rows written", got)
	}
}

// The check is only worth anything if the registry it reads is populated in
// the binary that serves requests, and it is filled by init() in packages
// cmd/api does not import - cmd/migrate blank-imports them, and cmd wires both
// subcommands into one binary.
//
// This cannot be asserted from an ordinary test: importing the version package
// to look at the registry would put it in the test binary's dependency graph
// and pass whatever the real binary links. So ask the build instead.
//
// Without this, dropping those blank imports leaves a check that reports
// "nothing pending" for every database forever, and every test above still
// passes.
func TestTheServingBinaryLinksTheMigrationRegistry(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", "go-admin").Output()
	if err != nil {
		t.Skipf("go list unavailable: %v", err)
	}
	deps := string(out)

	const versions = "go-admin/cmd/migrate/migration/version"
	if !strings.Contains(deps, versions+"\n") {
		t.Errorf("the main package does not link %s, so the schema check would "+
			"read an empty registry and report every database as current", versions)
	}

	// Negative control: a package the binary genuinely must not link, so that a
	// `deps` that somehow contained everything would fail here rather than pass
	// the assertion above for the wrong reason.
	const notLinked = "go-admin/tools/checksilent"
	if strings.Contains(deps, notLinked+"\n") {
		t.Errorf("%s is in the binary's dependency closure, so this test cannot "+
			"tell a real link from a query that matches anything", notLinked)
	}
}
