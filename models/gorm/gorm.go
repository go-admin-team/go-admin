package gorm

import (
	"gorm.io/gorm"

	"go-admin/models"
	"go-admin/models/tools"
)

func AutoMigrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate(new(models.CasbinRule),
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
		new(models.SysContent))
	if err != nil {
		return err
	}

	models.DataInit()
	return err
}

func CustomMigrate(db *gorm.DB, table interface{}) error {
	return db.AutoMigrate(&table)
}
