package system

import (
	"go-admin/common/models"
)

type SysDictData struct {
	models.ControlBy
	models.ModelTime

	DictCode  int    `json:"dictCode" gorm:"primaryKey;column:dict_code;autoIncrement;comment:主键编码"`
	DictSort  int    `json:"dictSort" gorm:"type:bigint(20);comment:DictSort"`
	DictLabel string `json:"dictLabel" gorm:"type:varchar(128);comment:DictLabel"`
	DictValue string `json:"dictValue" gorm:"type:varchar(255);comment:DictValue"`
	DictType  string `json:"dictType" gorm:"type:varchar(64);comment:DictType"`
	CssClass  string `json:"cssClass" gorm:"type:varchar(128);comment:CssClass"`
	ListClass string `json:"listClass" gorm:"type:varchar(128);comment:ListClass"`
	IsDefault string `json:"isDefault" gorm:"type:varchar(8);comment:IsDefault"`
	Status    string `json:"status" gorm:"type:varchar(4);comment:Status"`
	Default   string `json:"default" gorm:"type:varchar(8);comment:Default"`
	Remark    string `json:"remark" gorm:"type:varchar(255);comment:Remark"`
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
