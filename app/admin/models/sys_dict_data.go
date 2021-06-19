package models

import (
	"go-admin/common/models"
)

type SysDictData struct {
	DictCode  int    `json:"dictCode" gorm:"primaryKey;column:dict_code;autoIncrement;comment:主键编码"`
	DictSort  int    `json:"dictSort" gorm:"size:20;comment:DictSort"`
	DictLabel string `json:"dictLabel" gorm:"size:128;comment:DictLabel"`
	DictValue string `json:"dictValue" gorm:"size:255;comment:DictValue"`
	DictType  string `json:"dictType" gorm:"size:64;comment:DictType"`
	CssClass  string `json:"cssClass" gorm:"size:128;comment:CssClass"`
	ListClass string `json:"listClass" gorm:"size:128;comment:ListClass"`
	IsDefault string `json:"isDefault" gorm:"size:8;comment:IsDefault"`
	Status    int    `json:"status" gorm:"size:4;comment:Status"`
	Default   string `json:"default" gorm:"size:8;comment:Default"`
	Remark    string `json:"remark" gorm:"size:255;comment:Remark"`
	models.ControlBy
	models.ModelTime
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}

func (e *SysDictData) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDictData) GetId() interface{} {
	return e.DictCode
}
