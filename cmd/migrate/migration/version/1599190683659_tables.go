package version

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/models/tools"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fileName = filepath.Base(fileName)
	fileName = fileName[:len(fileName)-3]
	migration.Migrate.SetVersion(fileName, _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	err := db.Migrator().AutoMigrate(
		new(models.CasbinRule),
		new(models.SysDept),
		new(models.SysConfig),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(models.Menu),
		new(models.LoginLog),
		new(models.SysOperLog),
		new(models.RoleMenu),
		new(models.SysRoleDept),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Post),
		new(models.DictData),
		new(models.DictType),
		new(models.SysJob),
		new(models.SysConfig),
		new(models.SysSetting),
		new(models.SysFileDir),
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
