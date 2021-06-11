package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		new(models.CasbinRule),
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
		new(models.SysChinaAreaData),
		new(models.SysJob),
		new(models.SysConfig),
		new(models.SysSetting),
		new(models.SysFileDir),
		new(models.SysFileInfo),
		new(models.SysCategory),
		new(models.SysContent),
		new(models.SysApi),
		new(models.SysChinaAreaData),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}
