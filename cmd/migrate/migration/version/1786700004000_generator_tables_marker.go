package version

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Convert the two generator tables to the same delete marker every other table
// got in 1786700003000.
//
// They were left out of that list, and the mismatch is invisible until it is
// not: the tables are built from cmd/migrate/migration/models, whose ModelTime
// still carries a nullable gorm.DeletedAt, while the runtime models in
// app/other/models/tools embed common.ModelTime, which is the millisecond
// marker. GORM therefore queries them with deleted_at = 0 against a datetime
// column holding NULL, and every row is invisible - so the code generator
// listed no tables at all.
//
// tb_demo has the same shape and is deliberately not here: nothing reads it at
// runtime, so there is no mismatch to fix.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700004000GeneratorTablesMarker)
}

var generatorTables = []string{"sys_columns", "sys_tables"}

func _1786700004000GeneratorTablesMarker(db *gorm.DB, version string) error {
	for _, table := range generatorTables {
		if err := convertDeletedAt(db, table); err != nil {
			return fmt.Errorf("%s: %w", table, err)
		}
	}
	return db.Create(&common.Migration{Version: version}).Error
}
