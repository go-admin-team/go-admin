package models

import (
	"go-admin/common/models"
)

type SysChinaAreaData struct {
	models.Model

	PId        string    `json:"pId" gorm:"type:int(11);comment:上级编码"`
	Name       string    `json:"name" gorm:"type:varchar(128);comment:名称"`
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