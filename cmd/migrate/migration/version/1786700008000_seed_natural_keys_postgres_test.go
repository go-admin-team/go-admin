package version

import (
	"testing"
)

// postgresDB is defined in 1786700003000_soft_delete_marker_postgres_test.go.
//
// This file exists because refuseOnDuplicateApis's duplicate check is
// spelled with CONCAT(), a function this migration's design assumed
// PostgreSQL has carried since 9.1 but that nothing had run against a real
// PostgreSQL server before this test - only against the pure-Go SQLite
// driver, which happens to bundle a SQLite new enough to have grown its own
// CONCAT() only recently. A dialect where that assumption were wrong would
// otherwise only be discovered the first time an operator's install hit a
// genuine sys_api duplicate on PostgreSQL in production.
func TestSeedNaturalKeysRefusesDuplicateApisOnPostgres(t *testing.T) {
	db := postgresDB(t)
	t.Cleanup(func() { db.Migrator().DropTable(&oldSeedMenu{}, &oldSeedApi{}) })
	db.Migrator().DropTable(&oldSeedMenu{}, &oldSeedApi{})
	if err := db.AutoMigrate(&oldSeedMenu{}, &oldSeedApi{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	for i := 0; i < 2; i++ {
		if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
			t.Fatalf("seed duplicate %d: %v", i, err)
		}
	}

	err := seedNaturalKeys(db)
	if err == nil {
		t.Fatal("PostgreSQL accepted sys_api rows that already hold a duplicate (app_code, path, action)")
	}
	if !contains(err.Error(), "order") || !contains(err.Error(), "/api/v1/order") {
		t.Errorf("the error does not name the offending row: %v", err)
	}
	if db.Migrator().HasIndex(&oldSeedApi{}, "uk_sys_api_app_path_action_del") {
		t.Error("the unique index was built despite the migration refusing")
	}
}

// The success path, on the same server: both columns and both unique
// indexes have to actually build on PostgreSQL, not merely fail to error
// out on SQLite. Mirrors TestSeedNaturalKeysIsRepeatable's SQLite coverage.
func TestSeedNaturalKeysBuildsOnPostgres(t *testing.T) {
	db := postgresDB(t)
	t.Cleanup(func() { db.Migrator().DropTable(&oldSeedMenu{}, &oldSeedApi{}) })
	db.Migrator().DropTable(&oldSeedMenu{}, &oldSeedApi{})
	if err := db.AutoMigrate(&oldSeedMenu{}, &oldSeedApi{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	for i := 0; i < 2; i++ {
		if err := seedNaturalKeys(db); err != nil {
			t.Fatalf("migrate %d: %v", i, err)
		}
	}

	if !db.Migrator().HasColumn(&oldSeedMenu{}, "seed_code") {
		t.Error("sys_menu.seed_code was not added on PostgreSQL")
	}
	if !db.Migrator().HasIndex(&oldSeedMenu{}, "uk_sys_menu_app_seed_code_del") {
		t.Error("the sys_menu unique index was not built on PostgreSQL")
	}
	if !db.Migrator().HasIndex(&oldSeedApi{}, "uk_sys_api_app_path_action_del") {
		t.Error("the sys_api unique index was not built on PostgreSQL")
	}
}
