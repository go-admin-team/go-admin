package version

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	common "go-admin/common/models"
)

// oldSeedMenu/oldSeedApi are the shape of sys_menu/sys_api immediately
// before this migration: post-1786700003000 (deleted_at is the NOT NULL
// millisecond marker) and post-1786700006000 (app_code exists), but before
// seed_code or either unique index. They stand in for the real runtime
// models, which by the time this file is read already carry the columns
// this migration adds - the same relationship oldUser bears to sys_user in
// 1786700003000_soft_delete_marker_test.go.
type oldSeedMenu struct {
	MenuId    int    `gorm:"column:menu_id;primaryKey;autoIncrement"`
	AppCode   string `gorm:"column:app_code;type:varchar(64);not null;default:''"`
	DeletedAt int64  `gorm:"column:deleted_at;not null;default:0"`
}

func (oldSeedMenu) TableName() string { return "sys_menu" }

type oldSeedApi struct {
	Id        int    `gorm:"column:id;primaryKey;autoIncrement"`
	AppCode   string `gorm:"column:app_code;type:varchar(64);not null;default:''"`
	Path      string `gorm:"column:path;type:varchar(128)"`
	Action    string `gorm:"column:action;type:varchar(16)"`
	DeletedAt int64  `gorm:"column:deleted_at;not null;default:0"`
}

func (oldSeedApi) TableName() string { return "sys_api" }

func openSeedNaturalKeysDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&oldSeedMenu{}, &oldSeedApi{}, &common.Migration{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

// The host's own hand-placed menus, and every app-seeded row written
// before this column existed, have no seed_code at all - an unbounded
// number of those must coexist under the same app_code without tripping
// the new unique index (design doc §1.6: "NULL never treated as equal to
// NULL").
func TestSeedNaturalKeysToleratesManyPreExistingMenusWithNoSeedCode(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	for i := 0; i < 3; i++ {
		if err := db.Create(&oldSeedMenu{AppCode: ""}).Error; err != nil {
			t.Fatalf("seed pre-existing menu %d: %v", i, err)
		}
	}

	if err := seedNaturalKeys(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if !db.Migrator().HasColumn(&adminmodels.SysMenu{}, "SeedCode") {
		t.Fatal("sys_menu.seed_code was not added")
	}
}

// The point of adding seed_code at all: a second row with the same
// (app_code, seed_code) while both are live is what seedMenuTree's
// idempotency check depends on the database to reject if the Go-level
// check above it is ever bypassed or raced.
func TestSeedNaturalKeysMenuUniqueIndexBindsLiveRowsOnly(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	if err := seedNaturalKeys(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := db.Exec(
		"INSERT INTO sys_menu (app_code, seed_code, deleted_at) VALUES ('order', 'dir', 0)",
	).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	t.Run("a second live row with the same natural key is rejected", func(t *testing.T) {
		err := db.Exec(
			"INSERT INTO sys_menu (app_code, seed_code, deleted_at) VALUES ('order', 'dir', 0)",
		).Error
		if err == nil {
			t.Fatal("a duplicate (app_code, seed_code) was accepted while both rows were live")
		}
	})

	t.Run("the key is free again once the row is soft-deleted", func(t *testing.T) {
		if err := db.Exec("UPDATE sys_menu SET deleted_at = ? WHERE seed_code = 'dir'", time.Now().UnixMilli()).Error; err != nil {
			t.Fatalf("soft-delete: %v", err)
		}
		if err := db.Exec(
			"INSERT INTO sys_menu (app_code, seed_code, deleted_at) VALUES ('order', 'dir', 0)",
		).Error; err != nil {
			t.Errorf("the key stayed taken after its row was soft-deleted: %v", err)
		}
	})
}

// The demo site has already proven sys_api can hold historical duplicates;
// the migration has to name them and refuse, not let CREATE UNIQUE INDEX
// fail on an operator with no idea which rows to reconcile.
func TestSeedNaturalKeysRefusesDuplicateApis(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	for i := 0; i < 2; i++ {
		if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
			t.Fatalf("seed duplicate %d: %v", i, err)
		}
	}

	err := seedNaturalKeys(db)
	if err == nil {
		t.Fatal("the migration accepted sys_api rows that already hold a duplicate (app_code, path, action)")
	}
	if !contains(err.Error(), "order") || !contains(err.Error(), "/api/v1/order") {
		t.Errorf("the error does not name the offending row: %v", err)
	}
	if db.Migrator().HasIndex(&adminmodels.SysApi{}, "uk_sys_api_app_path_action_del") {
		t.Error("the unique index was built despite the migration refusing")
	}
	// sys_menu's column and index are independent of sys_api's outcome and
	// should already be in place - a partial failure here still leaves a
	// record of what succeeded, same as any other non-transactional DDL
	// migration in this package.
	if !db.Migrator().HasColumn(&adminmodels.SysMenu{}, "SeedCode") {
		t.Error("sys_menu.seed_code was not added even though only the sys_api step failed")
	}
}

// Only live rows count towards the duplicate check: a row a prior,
// unrelated soft-delete already retired does not block the index it will
// never occupy a slot in.
func TestSeedNaturalKeysIgnoresSoftDeletedApiDuplicates(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
		t.Fatalf("seed live row: %v", err)
	}
	if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET", DeletedAt: time.Now().UnixMilli()}).Error; err != nil {
		t.Fatalf("seed soft-deleted row: %v", err)
	}

	if err := seedNaturalKeys(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if !db.Migrator().HasIndex(&adminmodels.SysApi{}, "uk_sys_api_app_path_action_del") {
		t.Error("the unique index was not built")
	}
}

// The point of the sys_api index, mirroring
// TestSeedNaturalKeysMenuUniqueIndexBindsLiveRowsOnly above: a second live
// row is rejected, and the key is free again once the row is
// soft-deleted.
func TestSeedNaturalKeysApiUniqueIndexBindsLiveRowsOnly(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := seedNaturalKeys(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	t.Run("a second live row with the same natural key is rejected", func(t *testing.T) {
		err := db.Exec(
			"INSERT INTO sys_api (app_code, path, action, deleted_at) VALUES ('order', '/api/v1/order', 'GET', 0)",
		).Error
		if err == nil {
			t.Fatal("a duplicate (app_code, path, action) was accepted while both rows were live")
		}
	})

	t.Run("the key is free again once the row is soft-deleted", func(t *testing.T) {
		if err := db.Exec(
			"UPDATE sys_api SET deleted_at = ? WHERE path = '/api/v1/order'", time.Now().UnixMilli(),
		).Error; err != nil {
			t.Fatalf("soft-delete: %v", err)
		}
		if err := db.Exec(
			"INSERT INTO sys_api (app_code, path, action, deleted_at) VALUES ('order', '/api/v1/order', 'GET', 0)",
		).Error; err != nil {
			t.Errorf("the key stayed taken after its row was soft-deleted: %v", err)
		}
	})
}

// Running it twice must be safe: DDL does not roll back on MySQL, so an
// operator whose first attempt failed partway through (say, sys_menu's step
// succeeded and sys_api's refused) has nothing to do but run the whole
// migration again once the duplicates are reconciled.
func TestSeedNaturalKeysIsRepeatable(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	if err := db.Create(&oldSeedApi{AppCode: "order", Path: "/api/v1/order", Action: "GET"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := seedNaturalKeys(db); err != nil {
			t.Fatalf("migrate %d: %v", i, err)
		}
	}
}

// The wrapper's contract with Migrate.run(): the version is only recorded
// once the whole thing - both columns, both indexes - succeeded.
func TestSeedNaturalKeysWrapperRecordsTheVersion(t *testing.T) {
	db := openSeedNaturalKeysDB(t)
	if err := _1786700008000SeedNaturalKeys(db, "1786700008000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	var applied common.Migration
	if err := db.Where("version = ?", "1786700008000").First(&applied).Error; err != nil {
		t.Fatalf("sys_migration was not recorded: %v", err)
	}
}
