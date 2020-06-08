package gorm

import (
	"github.com/jinzhu/gorm"
	"go-admin/models"
	"go-admin/models/tools"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		new(models.CasbinRule),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(models.Dept),
		new(models.Menu),
		new(models.LoginLog),
		new(models.SysOperLog),
		new(models.RoleMenu),
		new(models.SysRoleDept),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Post),
		new(models.DictData),
		new(models.SysConfig),
		new(models.DictType),
	).Error
}
