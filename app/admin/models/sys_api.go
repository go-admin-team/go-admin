package models

import (
	// "gorm.io/gorm"

	"go-admin/common/models"
)

type SysApi struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Handle   string `json:"handle" gorm:"type:varchar(128);comment:handle"`
	Title    string `json:"title" gorm:"type:varchar(128);comment:标题"`
	Path     string `json:"path" gorm:"type:varchar(128);comment:地址"`
	Paths    string `json:"paths" gorm:"type:varchar(128);comment:Paths"`
	Action   string `json:"action" gorm:"type:varchar(16);comment:类型"`
	ParentId int    `json:"parentId" gorm:"comment:按钮id"`
	Sort     int    `json:"sort" gorm:"type:tinyint(4);comment:排序"`
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