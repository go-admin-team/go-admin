package version

import (
	"runtime"

	"github.com/go-admin-team/go-admin-core/sdk/config"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if config.DatabaseConfig.Driver == "mysql" {
			tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		err := tx.Migrator().AutoMigrate(
			new(models.SysDept),
			new(models.SysConfig),
			new(models.SysTables),
			new(models.SysColumns),
			new(models.SysMenu),
			new(models.SysLoginLog),
			new(models.SysOperaLog),
			new(models.SysRoleDept),
			new(models.SysUser),
			new(models.SysRole),
			new(models.SysPost),
			new(models.DictData),
			new(models.DictType),
			new(models.SysJob),
			new(models.SysApi),
			new(models.TbDemo),
		)
		if err != nil {
			return err
		}
		if err := models.InitDb(tx); err != nil {
			return err
		}
		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
