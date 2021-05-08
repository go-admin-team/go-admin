package models

import (
	"go-admin/common/models"
)

type SysCategory struct {
	models.Model
	Name   string `json:"name" gorm:"size:255;comment:名称"`   //
	Img    string `json:"img" gorm:"size:255;comment:图标"`    //
	Sort   int    `json:"sort" gorm:"size:4;comment:排序"`     //
	Status int    `json:"status" gorm:"size:1;comment:状态"`   //
	Remark string `json:"remark" gorm:"size:255;comment:备注"` //
	models.ControlBy
	models.ModelTime
}

func (SysCategory) TableName() string {
	return "sys_category"
}

func (e *SysCategory) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysCategory) GetId() interface{} {
	return e.Id
}
