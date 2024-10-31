package version_local

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"
	"runtime"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1653638869132Test)
}

/*
*
开发者项目的迁移脚本放在这个目录里，init写法参考version目录里的migrate或者自动生成
*/
func _1653638869132Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if config.DatabaseConfig.Driver == "mysql" {
			tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		err := tx.Migrator().AutoMigrate(
			new(models.TApiZl),
		)
		if err != nil {
			return err
		}
		// if err := models.InitDb(tx); err != nil {
		// 	return err
		// }
		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
