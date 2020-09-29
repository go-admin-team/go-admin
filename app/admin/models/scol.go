package models

import (
    "gorm.io/gorm"

	"go-admin/common/models"
)

type Scol struct {
    gorm.Model
    models.ControlBy
    
    Name string `json:"name" gorm:"type:varchar(255);comment:名称"` // 
}

func (Scol) TableName() string {
    return "scol"
}

func (e *Scol) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Scol) GetId() interface{} {
	return e.ID
}
