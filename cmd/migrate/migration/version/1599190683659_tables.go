package version

import (
	"runtime"

	"go-admin/app/admin/models/system"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/tools"
	"go-admin/cmd/migrate/migration"
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
		new(system.SysConfig),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(models.Menu),
		new(system.SysLoginLog),
		new(system.SysOperaLog),
		new(models.RoleMenu),
		new(models.SysRoleDept),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Post),
		new(models.DictData),
		new(models.DictType),
		new(models.SysJob),
		new(system.SysConfig),
		new(models.SysSetting),
		new(models.SysFileDir),
		new(models.SysFileInfo),
		new(models.SysCategory),
		new(models.SysContent),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}
