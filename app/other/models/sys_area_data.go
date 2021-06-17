package models

import (
	"go-admin/common/models"
)

type SysAreaData struct {
	Id       int                `json:"id" gorm:"primaryKey;comment:主键编码"`
	PId      int                `json:"pId" gorm:"size:11;comment:上级编码"`
	Name     string             `json:"name" gorm:"size:128;comment:名称"`
	Children []SysAreaData `json:"children,omitempty" gorm:"-"`
	models.ControlBy
	models.ModelTime
}

func (SysAreaData) TableName() string {
	return "sys_china_area_data"
}

func (e *SysAreaData) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysAreaData) GetId() interface{} {
	return e.Id
}
