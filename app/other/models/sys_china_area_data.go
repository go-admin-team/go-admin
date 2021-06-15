package models

import (
	"go-admin/common/models"
)

type SysChinaAreaData struct {
	Id       int                `json:"id" gorm:"primaryKey;comment:主键编码"`
	PId      int                `json:"pId" gorm:"size:11;comment:上级编码"`
	Name     string             `json:"name" gorm:"size:128;comment:名称"`
	Children []SysChinaAreaData `json:"children,omitempty" gorm:"-"`
	models.ControlBy
	models.ModelTime
}

func (SysChinaAreaData) TableName() string {
	return "sys_china_area_data"
}

func (e *SysChinaAreaData) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysChinaAreaData) GetId() interface{} {
	return e.Id
}
