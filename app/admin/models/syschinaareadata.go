package models

import (
	"go-admin/common/models"
	"time"
)

type SysChinaAreaData struct {
	models.Model
	models.ControlBy

	PId        string    `json:"pId" gorm:"type:int(11);comment:上级编码"`
	Name       string    `json:"name" gorm:"type:varchar(128);comment:名称"`
	CreateTime time.Time `json:"createTime" gorm:"type:timestamp;comment:CreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"type:timestamp;comment:UpdateTime"`
	DeleteTime time.Time `json:"deleteTime" gorm:"type:timestamp;comment:DeleteTime"`
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
