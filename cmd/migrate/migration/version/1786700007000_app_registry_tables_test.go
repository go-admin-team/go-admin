package version

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	common "go-admin/common/models"
)

func openAppRegistryDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&common.Migration{}); err != nil {
		t.Fatalf("automigrate sys_migration: %v", err)
	}
	return db
}

// The migration has to build both tables and record itself as applied -
// F2/F6's acceptance case is a row landing in either one, and neither is
// possible if the table it belongs to was never created.
func TestAppRegistryTablesAreCreated(t *testing.T) {
	db := openAppRegistryDB(t)

	if err := _1786700007000AppRegistryTables(db, "1786700007000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if !db.Migrator().HasTable(&adminmodels.SysApp{}) {
		t.Fatal("sys_app was not created")
	}
	if !db.Migrator().HasTable(&adminmodels.SysAppCasbinGrant{}) {
		t.Fatal("sys_app_casbin_grant was not created")
	}

	// A row that exercises every column, not just HasTable/HasColumn -
	// AutoMigrate can build a column with the wrong type and still report
	// that it exists.
	if err := db.Create(&adminmodels.SysApp{
		AppCode: "order", Name: "Order", Version: "v1", Description: "d", Author: "a",
		Requires: "payment", Pricing: "free", License: "MIT", Status: 1,
	}).Error; err != nil {
		t.Fatalf("insert sys_app: %v", err)
	}
	if err := db.Create(&adminmodels.SysAppCasbinGrant{
		AppCode: "order", Ptype: "p", V0: "admin", V1: "/api/v1/order", V2: "GET",
	}).Error; err != nil {
		t.Fatalf("insert sys_app_casbin_grant: %v", err)
	}

	var applied common.Migration
	if err := db.Where("version = ?", "1786700007000").First(&applied).Error; err != nil {
		t.Fatalf("sys_migration was not recorded: %v", err)
	}
}

// Running it twice must be safe: DDL does not roll back on MySQL, so an
// operator whose first attempt failed partway through has nothing to do but
// run it again. This calls AutoMigrate directly rather than the wrapper,
// which also inserts a sys_migration row that a second call would collide
// on - a collision Migrate.run() itself prevents by never calling a
// function twice for the same recorded version, so it is not this
// migration's job to tolerate.
func TestAppRegistryTablesAutoMigrateIsRepeatable(t *testing.T) {
	db := openAppRegistryDB(t)

	for i := 0; i < 3; i++ {
		if err := db.Migrator().AutoMigrate(
			new(adminmodels.SysApp),
			new(adminmodels.SysAppCasbinGrant),
		); err != nil {
			t.Fatalf("automigrate %d: %v", i, err)
		}
	}
}

// sys_app.app_code is the unique key G2 ("is app X installed") answers with
// - a second row for the same app code must be rejected, not tolerated.
func TestSysAppAppCodeIsUnique(t *testing.T) {
	db := openAppRegistryDB(t)
	if err := _1786700007000AppRegistryTables(db, "1786700007000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := db.Create(&adminmodels.SysApp{AppCode: "order", Name: "Order", Version: "v1"}).Error; err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if err := db.Create(&adminmodels.SysApp{AppCode: "order", Name: "Order dup", Version: "v1"}).Error; err == nil {
		t.Fatal("a second sys_app row with the same app_code was accepted")
	}
}

// sys_app_casbin_grant's unique index mirrors casbin_rule's own natural key
// (ptype,v0..v5) exactly - see design doc §3. A duplicate grant for the
// same rule must be rejected the same way gorm-adapter's own unique index
// on casbin_rule would reject it.
func TestSysAppCasbinGrantNaturalKeyIsUnique(t *testing.T) {
	db := openAppRegistryDB(t)
	if err := _1786700007000AppRegistryTables(db, "1786700007000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	grant := adminmodels.SysAppCasbinGrant{AppCode: "order", Ptype: "p", V0: "admin", V1: "/api/v1/order", V2: "GET"}
	if err := db.Create(&grant).Error; err != nil {
		t.Fatalf("first insert: %v", err)
	}
	dup := grant
	dup.Id = 0
	if err := db.Create(&dup).Error; err == nil {
		t.Fatal("a second sys_app_casbin_grant row with the same natural key was accepted")
	}

	// A grant for a different app, but the identical casbin natural key, is
	// exactly the collision two applications granting the same api/role
	// pair would produce - the natural key has to be the one thing that
	// rejects it, app_code is descriptive only and not part of the index.
	other := grant
	other.Id = 0
	other.AppCode = "another-app"
	if err := db.Create(&other).Error; err == nil {
		t.Fatal("a duplicate natural key under a different app_code was accepted")
	}
}
