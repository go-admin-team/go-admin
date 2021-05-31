package models

import (
	"go-admin/common/models"
)

type SysFileInfo struct {
	models.Model
	Type    string `json:"type" gorm:"size:255;comment:类型"`      //
	Name    string `json:"name" gorm:"size:255;comment:名称"`      //
	Size    string `json:"size" gorm:"size:11;comment:大小"`       //
	PId     int    `json:"p_id" gorm:"size:11;comment:目录"`       //
	Source  string `json:"source" gorm:"size:255;comment:来源"`    //
	Url     string `json:"url" gorm:"size:255;comment:地址"`       //
	FullUrl string `json:"full_url" gorm:"size:255;comment:全地址"` //
	models.ControlBy
	models.ModelTime
}

func (SysFileInfo) TableName() string {
	return "sys_file_info"
}

func (e *SysFileInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysFileInfo) GetId() interface{} {
	return e.Id
}
