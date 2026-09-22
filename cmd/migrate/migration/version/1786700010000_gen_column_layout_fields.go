package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/other/models/tools"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Add sys_columns.col_width and sys_columns.default_value for PRD 010 F1/F2
// (代码生成器前端模板迁移 Vue 3).
//
// col_width backs R2's column-width inference fallback and default_value
// backs R1/A6's "unconfigured rows still generate a usable page" guarantee -
// see docs-prd/010-代码生成器前端模板迁移Vue3/数据库变更.md §1.1 for why both
// defaults are sentinels (0 / "") rather than NULL: a non-pointer Go int/
// string field can never read NULL back out, and NULL would give
// "unconfigured" two representations instead of one.
//
// Ordered after 1786700003000, so this reads tools.SysColumns (the runtime
// model sys_columns's Update/GetPage/GetSysTablesInfo actually query through)
// rather than cmd/migrate/migration/models, matching every migration in this
// directory since sys_columns was converted - see
// 1786700004000_generator_tables_marker.go and schema_coverage_test.go's
// TestPostConversionMigrationsAvoidFrozenSeedModels.
//
// Hard prerequisite: tools.SysColumns must already declare ColWidth and
// DefaultValue (with the gorm tags in the doc above) by the time this file
// is compiled - AddColumn reads the column definition off the struct's own
// tag, not off anything in this file. Landing this migration without that
// model change first makes HasColumn/AddColumn silently do nothing (the
// field lookup fails and AddColumn returns an error naming the missing
// field), which fails loudly rather than silently - see the "no such field"
// error - so this is caught at migrate time, not left for a report later.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700010000GenColumnLayoutFields)
}

func _1786700010000GenColumnLayoutFields(db *gorm.DB, version string) error {
	m := db.Migrator()
	if !m.HasColumn(&tools.SysColumns{}, "ColWidth") {
		if err := m.AddColumn(&tools.SysColumns{}, "ColWidth"); err != nil {
			return err
		}
	}
	if !m.HasColumn(&tools.SysColumns{}, "DefaultValue") {
		if err := m.AddColumn(&tools.SysColumns{}, "DefaultValue"); err != nil {
			return err
		}
	}
	return db.Create(&common.Migration{Version: version}).Error
}
