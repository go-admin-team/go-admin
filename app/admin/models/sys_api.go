package models

import (
	// "gorm.io/gorm"

	"go-admin/common/models"
)

type SysApi struct {
	models.Model

	Name     string `json:"name" gorm:"type:varchar(128);comment:名称"`
	Title    string `json:"title" gorm:"type:varchar(128);comment:标题"`
	Path     string `json:"path" gorm:"type:varchar(128);comment:地址"`
	Paths    string `json:"paths" gorm:"type:varchar(128);comment:Paths"`
	Action   string `json:"action" gorm:"type:varchar(16);comment:类型"`
	ParentId string `json:"parentId" gorm:"type:smallint(6);comment:按钮id"`
	Sort     string `json:"sort" gorm:"type:tinyint(4);comment:排序"`
	models.ModelTime
	models.ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}
