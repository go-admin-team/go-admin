package models

import (
	// "gorm.io/gorm"

	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type SysApi struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Name     string `json:"name" gorm:"type:varchar(128);comment:名称"`
	Title    string `json:"title" gorm:"type:varchar(128);comment:标题"`
	Path     string `json:"path" gorm:"type:varchar(128);comment:地址"`
	Paths    string `json:"paths" gorm:"type:varchar(128);comment:Paths"`
	Action   string `json:"action" gorm:"type:varchar(16);comment:类型"`
	ParentId int    `json:"parentId" gorm:"type:smallint(6);comment:按钮id"`
	Sort     int    `json:"sort" gorm:"type:tinyint(4);comment:排序"`
	models.ModelTime
	models.ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}

func (e *SysApi) Create(tx *gorm.DB) (id int, err error) {
	result := tx.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err = result.Error
		return
	}
	err = e.InitPaths(tx)
	if err != nil {
		return
	}
	id = e.Id
	return
}

func (e *SysApi) InitPaths(tx *gorm.DB) (err error) {
	parentMenu := new(Menu)
	if e.ParentId != 0 {
		tx.Table("sys_api").Where("id = ?", e.ParentId).First(parentMenu)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return
		}
		e.Paths = parentMenu.Paths + "/" + pkg.IntToString(e.Id)
	} else {
		e.Paths = "/0/" + pkg.IntToString(e.Id)
	}
	tx.Table("sys_api").Where("id = ?", e.ParentId).Update("paths", e.Paths)
	return
}
