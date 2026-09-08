package version

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Give seed.SeedMenus's two write paths (seedApis, seedMenuTree in
// app/admin/service/seed.go) a real natural key to check before inserting,
// so a retried, partially-failed migration (see the design doc
// docs-prd/008-应用清单与安装器/数据库变更.md §1.5/§1.6) does not insert the
// same row twice. This has already happened in production once (duplicate
// sys_menu/casbin_rule rows on the demo site), not a theoretical risk.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700008000SeedNaturalKeys)
}

func _1786700008000SeedNaturalKeys(db *gorm.DB, version string) error {
	if err := seedNaturalKeys(db); err != nil {
		return err
	}
	return db.Create(&common.Migration{Version: version}).Error
}

// seedNaturalKeys is split out from the wrapper above so tests can call it
// against a database that only has sys_menu/sys_api, without also standing
// up sys_migration - and so it can be called more than once in the same
// test to prove the re-run tolerance the doc comment above promises: DDL
// does not roll back on MySQL, so an operator whose first attempt failed
// partway through has nothing to do but run the whole migration again.
func seedNaturalKeys(db *gorm.DB) error {
	m := db.Migrator()

	// sys_menu.seed_code is a brand-new column: every existing row becomes
	// NULL, and NULL never collides in the unique index built below, so
	// this needs no pre-check.
	if !m.HasColumn(&adminmodels.SysMenu{}, "SeedCode") {
		if err := m.AddColumn(&adminmodels.SysMenu{}, "SeedCode"); err != nil {
			return err
		}
	}
	if !m.HasIndex(&adminmodels.SysMenu{}, "uk_sys_menu_app_seed_code_del") {
		if err := db.Exec(
			"CREATE UNIQUE INDEX uk_sys_menu_app_seed_code_del ON sys_menu (app_code, seed_code, deleted_at)",
		).Error; err != nil {
			return err
		}
	}

	// sys_api reuses existing, already-populated columns, which the demo
	// site has already proven can hold duplicates. Refuse rather than let
	// CREATE UNIQUE INDEX fail on an operator with no idea which rows to
	// reconcile - same shape as 1786700003000_soft_delete_marker.go's
	// refuseOnDuplicates.
	if err := refuseOnDuplicateApis(db); err != nil {
		return err
	}
	if !m.HasIndex(&adminmodels.SysApi{}, "uk_sys_api_app_path_action_del") {
		if err := db.Exec(
			"CREATE UNIQUE INDEX uk_sys_api_app_path_action_del ON sys_api (app_code, path, action, deleted_at)",
		).Error; err != nil {
			return err
		}
	}

	return nil
}

// refuseOnDuplicateApis reports the (app_code, path, action) values that
// would make the unique index impossible, rather than the index failing to
// build and saying only that it did. Only live rows count: a soft-deleted
// duplicate does not block the index it will never occupy a slot in.
func refuseOnDuplicateApis(db *gorm.DB) error {
	var dupes []string
	if err := db.Raw(
		`SELECT CONCAT(app_code, '|', path, '|', action) FROM sys_api
		 WHERE deleted_at = 0 GROUP BY app_code, path, action HAVING COUNT(*) > 1`,
	).Scan(&dupes).Error; err != nil {
		return fmt.Errorf("checking sys_api for duplicates: %w", err)
	}
	if len(dupes) > 0 {
		return fmt.Errorf(
			"sys_api already holds duplicate (app_code,path,action) %v; reconcile them before this migration can add its unique index",
			dupes)
	}
	return nil
}
