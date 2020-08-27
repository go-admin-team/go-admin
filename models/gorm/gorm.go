package gorm

import (
	"go-admin/models"
	"go-admin/models/tools"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(new(models.CasbinRule))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysDept))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysConfig))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(tools.SysTables))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(tools.SysColumns))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Menu))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.LoginLog))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysOperLog))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.RoleMenu))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysRoleDept))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysUser))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysRole))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Post))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.DictData))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.DictType))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysJob))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysConfig))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysSetting))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysFileDir))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysCategory))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysContent))
	if err != nil {
		return err
	}

	models.DataInit()
	return err
}

func CustomMigrate(db *gorm.DB, table interface{}) error {
	return db.AutoMigrate(&table)
}
