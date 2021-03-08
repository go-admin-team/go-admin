package models

import (
   // "gorm.io/gorm"

	"go-admin/common/models"
)

type WfProcessClassify struct {
    models.Model
    models.ControlBy
    
    Name string `json:"name" gorm:"type:varchar(255);comment:名称"` 
}

func (WfProcessClassify) TableName() string {
    return "wf_process_classify"
}

func (e *WfProcessClassify) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *WfProcessClassify) GetId() interface{} {
	return e.Id
}