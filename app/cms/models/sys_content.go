package models

import (
	"go-admin/common/models"
)

type SysContent struct {
	models.Model
	CateId  int    `json:"cateId" gorm:"type:int(11);comment:分类id"`
	Name    string `json:"name" gorm:"type:varchar(255);comment:名称"`
	Status  int    `json:"status" gorm:"type:int(1);comment:状态"`
	Img     string `json:"img" gorm:"type:varchar(255);comment:图片"`
	Content string `json:"content" gorm:"type:text;comment:内容"`
	Remark  string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Sort    int    `json:"sort" gorm:"type:int(4);comment:排序"`
	models.ControlBy
	models.ModelTime
}

// TableName
func (SysContent) TableName() string {
	return "sys_content"
}

// Generate
func (e *SysContent) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysContent) GetId() interface{} {
	return e.Id
}
