package models

import (
    "gorm.io/gorm"

	"go-admin/common/models"
)

type SysCategory struct {
    gorm.Model
    models.ControlBy
    
    Name string `json:"name" gorm:"type:varchar(255);comment:名称"` // 
    Img string `json:"img" gorm:"type:varchar(255);comment:图标"` // 
    Sort string `json:"sort" gorm:"type:int(4);comment:排序"` // 
    Status string `json:"status" gorm:"type:int(1);comment:状态"` // 
    Remark string `json:"remark" gorm:"type:varchar(255);comment:备注"` // 
}

func (SysCategory) TableName() string {
    return "sys_category"
}

func (e *SysCategory) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysCategory) GetId() interface{} {
	return e.ID
}
