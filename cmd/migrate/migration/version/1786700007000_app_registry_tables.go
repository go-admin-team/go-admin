package version

import (
	"runtime"

	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Create sys_app (PRD 008 F2) and sys_app_casbin_grant (F4/F6's casbin
// attribution ledger - see the design doc's (docs-prd/008-应用清单与安装器/
// 数据库变更.md) §2.2/§3 for why casbin_rule itself is not touched:
// gorm-adapter's SavePolicyCtx truncates and reloads that table from its
// in-memory model, and any column this migration added to it would be
// silently zeroed the first time anything calls SavePolicy.
//
// Ordered after 1786700003000 (the soft-delete conversion), so importing
// cmd/migrate/migration/models is banned here - see
// schema_coverage_test.go's TestPostConversionMigrationsAvoidFrozenSeedModels.
// Both new tables are AutoMigrate'd from their runtime model shape under
// app/admin/models directly, which is also why neither one is added to
// 1786700003000's frozen softDeleteTables list: neither embeds
// common.ModelTime in the first place (see design doc §1.1).
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700007000AppRegistryTables)
}

func _1786700007000AppRegistryTables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Migrator().AutoMigrate(
			new(adminmodels.SysApp),
			new(adminmodels.SysAppCasbinGrant),
		); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}
