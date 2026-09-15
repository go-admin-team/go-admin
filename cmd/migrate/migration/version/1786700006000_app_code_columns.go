package version

import (
	"runtime"

	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Add sys_menu.app_code and sys_api.app_code ahead of PRD 006 F9's Seeder.
//
// Every row a third-party application's migration writes through
// seed.SeedMenus must be attributable to the app that wrote it, so
// installing, auditing, or removing one application does not require
// guessing which rows belong to it - see go-admin-core's docs/contract.md,
// "Application-supplied menu and API entries", for the requirement this
// satisfies.
//
// Ordered after 1786700003000, so importing cmd/migrate/migration/models is
// banned here (see schema_coverage_test.go's
// TestPostConversionMigrationsAvoidFrozenSeedModels): AddColumn instead
// reads the runtime models' own gorm tags directly, which is also what
// makes the column this adds match the one the admin Seeder writes through
// those same structs.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700006000AppCodeColumns)
}

func _1786700006000AppCodeColumns(db *gorm.DB, version string) error {
	m := db.Migrator()
	if !m.HasColumn(&adminmodels.SysMenu{}, "AppCode") {
		if err := m.AddColumn(&adminmodels.SysMenu{}, "AppCode"); err != nil {
			return err
		}
	}
	if !m.HasColumn(&adminmodels.SysApi{}, "AppCode") {
		if err := m.AddColumn(&adminmodels.SysApi{}, "AppCode"); err != nil {
			return err
		}
	}
	return db.Create(&common.Migration{Version: version}).Error
}
