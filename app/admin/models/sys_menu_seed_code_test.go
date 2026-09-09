package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ptr(s string) *string { return &s }

// uk_sys_menu_app_seed_code_del covers (app_code, seed_code, deleted_at) and
// is created by 1786700008000, not from a struct tag. It cannot come from a
// tag: a named uniqueIndex collects every field carrying that name, and
// deleted_at lives in the shared ModelTime embed that no single model can tag.
//
// Naming it on SeedCode alone anyway produced a unique index on seed_code by
// itself under the same name - stricter than the real one - and AutoMigrate
// here would create that one, after which the migration's HasIndex guard
// finds the name taken and leaves the wrong index in place. Nothing in
// production AutoMigrates this model (the initial table migration uses a
// frozen snapshot that has neither column), which is why this never showed up
// as a broken database; it showed up the first time a test built the schema
// from the live model and seeded two applications.
func TestSysMenuDeclaresNoSeedCodeIndexOfItsOwn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&SysMenu{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	if db.Migrator().HasIndex(&SysMenu{}, "uk_sys_menu_app_seed_code_del") {
		t.Error("AutoMigrate created uk_sys_menu_app_seed_code_del from a tag; " +
			"the migration's HasIndex guard will now skip the composite index it should create")
	}

	// Two applications, the same seed code. The real index allows it because
	// app_code is part of the key; an index on seed_code alone does not.
	for _, app := range []string{"order", "crm"} {
		row := SysMenu{MenuName: app + "Dir", AppCode: app, SeedCode: ptr("dir")}
		if err := db.Create(&row).Error; err != nil {
			t.Fatalf("%s could not use the seed code \"dir\": %v", app, err)
		}
	}
}
