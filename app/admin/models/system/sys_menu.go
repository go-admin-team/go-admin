package system

import "go-admin/common/models"

type SysMenu struct {
	MenuId     int       `json:"menuId" gorm:"primaryKey;autoIncrement"`
	MenuName   string    `json:"menuName" gorm:"size:128;"`
	Title      string    `json:"title" gorm:"size:128;"`
	Icon       string    `json:"icon" gorm:"size:128;"`
	Path       string    `json:"path" gorm:"size:128;"`
	Paths      string    `json:"paths" gorm:"size:128;"`
	MenuType   string    `json:"menuType" gorm:"size:1;"`
	Action     string    `json:"action" gorm:"size:16;"`
	Permission string    `json:"permission" gorm:"size:255;"`
	ParentId   int       `json:"parentId" gorm:"size:11;"`
	NoCache    bool      `json:"noCache" gorm:"size:8;"`
	Breadcrumb string    `json:"breadcrumb" gorm:"size:255;"`
	Component  string    `json:"component" gorm:"size:255;"`
	Sort       int       `json:"sort" gorm:"size:4;"`
	Visible    string    `json:"visible" gorm:"size:1;"`
	IsFrame    string    `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	DataScope  string    `json:"dataScope" gorm:"-"`
	Params     string    `json:"params" gorm:"-"`
	RoleId     int       `gorm:"-"`
	Children   []SysMenu `json:"children" gorm:"-"`
	IsSelect   bool      `json:"is_select" gorm:"-"`
	models.ControlBy
	models.ModelTime
}

type SysMenus struct {
	MenuId     int       `json:"menuId" gorm:"column:menu_id;primaryKey;autoIncrement;"`
	MenuName   string    `json:"menuName" gorm:"column:menu_name"`
	Title      string    `json:"title" gorm:"column:title"`
	Icon       string    `json:"icon" gorm:"column:icon"`
	Path       string    `json:"path" gorm:"column:path"`
	MenuType   string    `json:"menuType" gorm:"column:menu_type"`
	Action     string    `json:"action" gorm:"column:action"`
	Permission string    `json:"permission" gorm:"column:permission"`
	ParentId   int       `json:"parentId" gorm:"column:parent_id"`
	NoCache    bool      `json:"noCache" gorm:"column:no_cache"`
	Breadcrumb string    `json:"breadcrumb" gorm:"column:breadcrumb"`
	Component  string    `json:"component" gorm:"column:component"`
	Sort       int       `json:"sort" gorm:"column:sort"`
	Visible    string    `json:"visible" gorm:"column:visible"`
	Children   []SysMenu `json:"children" gorm:"-"`
	models.ControlBy
	models.ModelTime
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

func (e *SysMenu) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysMenu) GetId() interface{} {
	return e.MenuId
}
