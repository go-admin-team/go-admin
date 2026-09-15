package version

import (
	"testing"

	common "go-admin/common/models"

	adminmodels "go-admin/app/admin/models"
)

// postgresDB is defined in 1786700003000_soft_delete_marker_postgres_test.go
// and shared across this package's PostgreSQL-only tests.
//
// This migration is plain AutoMigrate on two brand-new tables, unlike
// 1786700003000's DROP INDEX (go-admin#919's actual defect), so there is no
// dialect-specific SQL here for AutoMigrate itself to get wrong on
// PostgreSQL specifically. What is worth a real PostgreSQL run is
// 1786700008000's CONCAT()-based duplicate check next door - PostgreSQL has
// had CONCAT() since 9.1, but it was never verified against a real server
// until this file, only inferred from documentation - and the same
// AutoMigrate call this test makes, so a schema/character-set mistake
// AutoMigrate might make silently on a dialect nobody ran it against here
// has somewhere to surface.
func TestAppRegistryTablesAreCreatedOnPostgres(t *testing.T) {
	db := postgresDB(t)
	const version = "1786700007000-pg"
	cleanup := func() {
		db.Migrator().DropTable(&adminmodels.SysAppCasbinGrant{}, &adminmodels.SysApp{})
		// Only this test's own row, not the whole shared sys_migration
		// table: postgresDB points at a real, persistent database (unlike
		// the SQLite tests' fresh in-memory one per run), so a version left
		// behind by a previous run of this same binary collides with the
		// wrapper's own INSERT the next time this test runs.
		db.Exec("DELETE FROM sys_migration WHERE version = ?", version)
	}
	t.Cleanup(cleanup)
	cleanup()
	if err := db.AutoMigrate(&common.Migration{}); err != nil {
		t.Fatalf("automigrate sys_migration: %v", err)
	}

	if err := _1786700007000AppRegistryTables(db, version); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := db.Create(&adminmodels.SysApp{AppCode: "order", Name: "Order", Version: "v1"}).Error; err != nil {
		t.Fatalf("insert sys_app: %v", err)
	}
	if err := db.Create(&adminmodels.SysApp{AppCode: "order", Name: "dup", Version: "v1"}).Error; err == nil {
		t.Fatal("a second sys_app row with the same app_code was accepted on PostgreSQL")
	}

	grant := adminmodels.SysAppCasbinGrant{AppCode: "order", Ptype: "p", V0: "admin", V1: "/api/v1/order", V2: "GET"}
	if err := db.Create(&grant).Error; err != nil {
		t.Fatalf("insert sys_app_casbin_grant: %v", err)
	}
	dup := grant
	dup.Id = 0
	if err := db.Create(&dup).Error; err == nil {
		t.Fatal("a second sys_app_casbin_grant row with the same natural key was accepted on PostgreSQL")
	}
}
