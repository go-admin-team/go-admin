package version

import (
	"os"
	"testing"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
)

// sqlserverDSNEnv points these tests at a database. They skip without it, so
// a developer with no SQL Server running still gets a green run.
//
// This file exists for the same reason the PostgreSQL one does, one driver
// further along. The rest of the package runs on SQLite, where the defect it
// covers cannot happen: SQLite, MySQL and PostgreSQL all treat two NULLs in a
// unique index as different values, and SQL Server treats them as equal and
// permits one. A suite that never pointed at SQL Server reported success for
// a migration that could not be applied to any SQL Server database at all,
// new or old.
const sqlserverDSNEnv = "GO_ADMIN_TEST_SQLSERVER_DSN"

func sqlserverDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv(sqlserverDSNEnv)
	if dsn == "" {
		// Skipping locally is the point; skipping in CI is the failure this
		// file exists to prevent.
		if os.Getenv("CI") != "" {
			t.Fatalf("%s is not set while CI is: the SQL Server migration tests must not skip here", sqlserverDSNEnv)
		}
		t.Skipf("%s is not set; skipping the SQL Server migration tests", sqlserverDSNEnv)
	}

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connecting to %s: %v", sqlserverDSNEnv, err)
	}
	return db
}

// freshSQLServerTables drops and rebuilds the two tables this migration
// touches, so a rerun does not inherit the previous run's index.
func freshSQLServerTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, m := range []any{&adminmodels.SysMenu{}, &adminmodels.SysApi{}} {
		if db.Migrator().HasTable(m) {
			if err := db.Migrator().DropTable(m); err != nil {
				t.Fatalf("dropping: %v", err)
			}
		}
	}
	if err := db.AutoMigrate(&adminmodels.SysMenu{}, &adminmodels.SysApi{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
}

// The migration completes on SQL Server.
//
// It did not. Five menus with no seed_code is what 1786700001000 leaves on
// every database, and the unfiltered index rejects the second of them:
//
//	Msg 1505 ... duplicate key ... The duplicate key value is (, <NULL>, 0).
func TestSeedNaturalKeysOnSQLServer(t *testing.T) {
	db := sqlserverDB(t)
	freshSQLServerTables(t, db)

	// Three rows in the state 1786700006000 leaves behind: an app_code that
	// defaulted to empty, no seed_code, and live.
	for _, name := range []string{"one", "two", "three"} {
		if err := db.Exec(
			"INSERT INTO sys_menu (menu_name, app_code, deleted_at) VALUES (?, '', 0)", name,
		).Error; err != nil {
			t.Fatalf("seeding %s: %v", name, err)
		}
	}
	// sys_api's key has two nullable columns and either one is enough to
	// collide, so both shapes are here. Two rows missing both, and two more
	// that have a path and no action: a filter naming only path would let
	// that second pair back into the index, where their equal NULLs collide.
	for i := 0; i < 2; i++ {
		if err := db.Exec("INSERT INTO sys_api (app_code, deleted_at) VALUES ('', 0)").Error; err != nil {
			t.Fatalf("seeding sys_api: %v", err)
		}
		if err := db.Exec(
			"INSERT INTO sys_api (app_code, path, deleted_at) VALUES ('', '/api/v1/half', 0)",
		).Error; err != nil {
			t.Fatalf("seeding a sys_api row with no action: %v", err)
		}
	}

	if err := seedNaturalKeys(db); err != nil {
		t.Fatalf("seedNaturalKeys on SQL Server: %v", err)
	}
	for _, name := range []string{"uk_sys_menu_app_seed_code_del", "uk_sys_api_app_path_action_del"} {
		var model any = &adminmodels.SysMenu{}
		if name == "uk_sys_api_app_path_action_del" {
			model = &adminmodels.SysApi{}
		}
		if !db.Migrator().HasIndex(model, name) {
			t.Errorf("%s was not created", name)
		}
	}

	// Rows that do carry a seed code still cannot collide - the filter takes
	// the ones with no value out of the index, it does not turn the index off.
	code := "dir"
	first := adminmodels.SysMenu{MenuName: "d1", AppCode: "order", SeedCode: &code}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("first seeded menu: %v", err)
	}
	second := adminmodels.SysMenu{MenuName: "d2", AppCode: "order", SeedCode: &code}
	if err := db.Create(&second).Error; err == nil {
		t.Error("a duplicate (app_code, seed_code) was accepted; the filtered index is not enforcing anything")
	}
	// A different app may reuse the same seed code, which is why the key is
	// composite in the first place.
	other := adminmodels.SysMenu{MenuName: "d3", AppCode: "crm", SeedCode: &code}
	if err := db.Create(&other).Error; err != nil {
		t.Errorf("another app could not reuse the seed code: %v", err)
	}
}

// The control. Without the filter the statement fails on this engine, so the
// test above is passing because of the fix rather than because SQL Server
// turned out not to mind.
func TestSQLServerRejectsTheUnfilteredIndex(t *testing.T) {
	db := sqlserverDB(t)
	freshSQLServerTables(t, db)

	for _, name := range []string{"one", "two"} {
		if err := db.Exec(
			"INSERT INTO sys_menu (menu_name, app_code, deleted_at) VALUES (?, '', 0)", name,
		).Error; err != nil {
			t.Fatalf("seeding %s: %v", name, err)
		}
	}

	err := db.Exec(uniqueIndexOverNullable("postgres",
		"uk_unfiltered_probe", "sys_menu", "app_code, seed_code, deleted_at", "seed_code")).Error
	if err == nil {
		t.Fatal("SQL Server accepted two NULLs in a unique index; the filter this migration adds is not needed")
	}
	t.Logf("as expected: %v", err)
}
