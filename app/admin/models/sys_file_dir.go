package models

import (
	"go-admin/common/models"
)

type SysFileDir struct {
	models.Model
	Label    string       `json:"label" gorm:"type:varchar(255);comment:目录名称"` // 目录名称
	PId      int          `json:"pId" gorm:"type:int(11);comment:上级目录"`        // 上级目录
	Sort     string       `json:"sort" gorm:"type:bigint(20);comment:排序"`      // 排序
	Path     string       `json:"path" gorm:"type:varchar(255);comment:路径"`    // 路径
	Children []SysFileDir `json:"children,omitempty" gorm:"-"`                 // 下级信息
	models.ControlBy
	models.ModelTime
}

type SysFileDirL struct {
	models.Model
	Label string `json:"label" gorm:"type:varchar(255);comment:目录名称"` // 目录名称
	PId   int    `json:"pId" gorm:"type:int(11);comment:上级目录"`        // 上级目录
	Sort  string `json:"sort" gorm:"type:bigint(20);comment:排序"`      // 排序
	Path  string `json:"path" gorm:"type:varchar(255);comment:路径"`    // 路径
	models.ControlBy
	models.ModelTime
	Children []SysFileDirL `json:"children,omitempty" gorm:"-"` // 下级信息
}

func (SysFileDir) TableName() string { /**/
	return "sys_file_dir"
}

func (e *SysFileDir) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysFileDir) GetId() interface{} {
	return e.Id
}
