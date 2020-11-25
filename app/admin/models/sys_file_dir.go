package models

import (
	"gorm.io/gorm"

	"go-admin/common/models"
)

type SysFileDir struct {
	gorm.Model
	models.ControlBy

	Label string `json:"label" gorm:"type:varchar(255);comment:目录名称"` // 目录名称
	PId   uint   `json:"pId" gorm:"type:int(11);comment:上级目录"`        // 上级目录
	//Sort  string `json:"sort" gorm:"type:bigint(20);comment:排序"`      // 排序
	Path  string `json:"path" gorm:"type:varchar(255);comment:路径"`    // 路径
}

type SysFileDirL struct {
	gorm.Model
	models.ControlBy

	Label    string        `json:"label" gorm:"type:varchar(255);comment:目录名称"` // 目录名称
	PId      uint          `json:"pId" gorm:"type:int(11);comment:上级目录"`        // 上级目录
	//Sort     string        `json:"sort" gorm:"type:bigint(20);comment:排序"`      // 排序
	Path     string        `json:"path" gorm:"type:varchar(255);comment:路径"`    // 路径
	Children []SysFileDirL `json:"children" gorm:"-"`                           // 下级信息
}

func (SysFileDir) TableName() string { /**/
	return "sys_file_dir"
}

func (e *SysFileDir) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysFileDir) GetId() interface{} {
	return e.ID
}
