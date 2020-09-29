package models

import (
    "gorm.io/gorm"
	"go-admin/common/models"
)

type Student struct {
    gorm.Model
    models.ControlBy
    
    Name string `json:"name" gorm:"type:varchar(255);comment:姓名"` // 
}

func (Student) TableName() string {
    return "student"
}

func (e *Student) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Student) GetId() interface{} {
	return e.ID
}
