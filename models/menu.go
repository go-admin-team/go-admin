package models

import (
	"errors"
	orm "go-admin/database"
	"go-admin/utils"
)

type Menu struct {
	MenuId     int64  `json:"menuId" gorm:"column:menu_id;primary_key"`
	MenuName   string `json:"menuName" gorm:"column:menu_name"`
	Title      string `json:"title" gorm:"column:title"`
	Icon       string `json:"icon" gorm:"column:icon"`
	Path       string `json:"path" gorm:"column:path"`
	Paths      string `json:"paths" gorm:"column:paths"`
	MenuType   string `json:"menuType" gorm:"column:menu_type"`
	Action     string `json:"action" gorm:"column:action"`
	Permission string `json:"permission" gorm:"column:permission"`
	ParentId   int64  `json:"parentId" gorm:"column:parent_id"`
	NoCache    bool   `json:"noCache" gorm:"column:no_cache"`
	Breadcrumb string `json:"breadcrumb" gorm:"column:breadcrumb"`
	Component  string `json:"component" gorm:"column:component"`
	Sort       int    `json:"sort" gorm:"column:sort"`

	Visible  string `json:"visible" gorm:"column:visible"`
	Children []Menu `json:"children" gorm:"-"`
	IsSelect bool   `json:"is_select" gorm:"-"`
	RoleId   int64  `gorm:"-"`

	CreateBy   string `json:"createBy" gorm:"column:create_by"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`
	UpdateTime string `json:"updateTime" gorm:"column:update_time"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"-"`
	IsDel      int `json:"isDel" gorm:"column:is_del"`
}
type MenuLable struct {
	Id       int64       `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []MenuLable `json:"children" gorm:"-"`
}

type Menus struct {
	MenuId     int64  `json:"menuId" gorm:"column:menu_id;primary_key"`
	MenuName   string `json:"menuName" gorm:"column:menu_name"`
	Title      string `json:"title" gorm:"column:title"`
	Icon       string `json:"icon" gorm:"column:icon"`
	Path       string `json:"path" gorm:"column:path"`
	MenuType   string `json:"menuType" gorm:"column:menu_type"`
	Action     string `json:"action" gorm:"column:action"`
	Permission string `json:"permission" gorm:"column:permission"`
	ParentId   int64  `json:"parentId" gorm:"column:parent_id"`
	NoCache    bool   `json:"noCache" gorm:"column:no_cache"`
	Breadcrumb string `json:"breadcrumb" gorm:"column:breadcrumb"`
	Component  string `json:"component" gorm:"column:component"`
	Sort       int    `json:"sort" gorm:"column:sort"`

	Visible  string `json:"visible" gorm:"column:visible"`
	Children []Menu `json:"children" gorm:"-"`

	CreateBy   string `json:"createBy" gorm:"column:create_by"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`
	UpdateTime string `json:"updateTime" gorm:"column:update_time"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"-"`
	IsDel      int `json:"isDel" gorm:"column:is_del"`
}

type MenuRole struct {
	Menus
	IsSelect bool `json:"is_select" gorm:"-"`
}

type MS []Menu

func (e *Menu) GetByMenuId() (Menu Menu, err error) {

	table := orm.Eloquent.Table("sys_menu")
	table = table.Where("is_del = ?", 0)
	table = table.Where("menu_id = ?", e.MenuId)
	if err = table.Find(&Menu).Error; err != nil {
		return
	}
	return
}

func (e *Menu) SetMenu() (m []Menu, err error) {
	menulist, err := e.GetPage()

	m = make([]Menu, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentId != 0 {
			continue
		}
		menusInfo := DiguiMenu(&menulist, menulist[i])

		m = append(m, menusInfo)
	}
	return
}

func DiguiMenu(menulist *[]Menu, menu Menu) Menu {
	list := *menulist

	min := make([]Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuId != list[j].ParentId {
			continue
		}
		mi := Menu{}
		mi.MenuId = list[j].MenuId
		mi.MenuName = list[j].MenuName
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.MenuType = list[j].MenuType
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentId = list[j].ParentId
		mi.NoCache = list[j].NoCache
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.CreateTime = list[j].CreateTime
		mi.UpdateTime = list[j].UpdateTime
		mi.IsDel = list[j].IsDel
		mi.Visible = list[j].Visible
		mi.Children = []Menu{}

		if mi.MenuType != "F" {
			ms := DiguiMenu(menulist, mi)
			min = append(min, ms)

		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}

func (e *Menu) SetMenuLable() (m []MenuLable, err error) {
	menulist, err := e.Get()

	m = make([]MenuLable, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentId != 0 {
			continue
		}
		e := MenuLable{}
		e.Id = menulist[i].MenuId
		e.Label = menulist[i].Title
		menusInfo := DiguiMenuLable(&menulist, e)

		m = append(m, menusInfo)
	}
	return
}

func DiguiMenuLable(menulist *[]Menu, menu MenuLable) MenuLable {
	list := *menulist

	min := make([]MenuLable, 0)
	for j := 0; j < len(list); j++ {

		if menu.Id != list[j].ParentId {
			continue
		}
		mi := MenuLable{}
		mi.Id = list[j].MenuId
		mi.Label = list[j].Title
		mi.Children = []MenuLable{}
		if list[j].MenuType != "F" {
			ms := DiguiMenuLable(menulist, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}

func (e *Menu) SetMenuRole(rolename string) (m []Menu, err error) {

	menulist, err := e.GetByRoleName(rolename)

	m = make([]Menu, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentId != 0 {
			continue
		}
		menusInfo := DiguiMenu(&menulist, menulist[i])

		m = append(m, menusInfo)
	}
	return
}

func (menu *MenuRole) Get() (Menus []MenuRole, err error) {
	table := orm.Eloquent.Table("sys_menu")
	table = table.Where("is_del = ?", 0)
	if menu.MenuName != "" {
		table = table.Where("menu_name = ?", menu.MenuName)
	}
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) GetByRoleName(rolename string) (Menus []Menu, err error) {
	table := orm.Eloquent.Table("sys_menu").Select("sys_menu.*").Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id")
	table = table.Where("is_del = ? and sys_role_menu.role_name=? and menu_type in ('M','C')", 0, rolename)
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) Get() (Menus []Menu, err error) {
	table := orm.Eloquent.Table("sys_menu")
	table = table.Where("is_del = ?", 0)
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.Path != "" {
		table = table.Where("path = ?", e.Path)
	}
	if e.Action != "" {
		table = table.Where("action = ?", e.Action)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}

	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) GetPage() (Menus []Menu, err error) {
	table := orm.Eloquent.Table("sys_menu")
	table = table.Where("is_del = ?", 0)
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}
	if e.Visible != "" {
		table = table.Where("visible = ?", e.Visible)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_menu", table)

	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) Create() (id int64, err error) {
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	e.IsDel = 0
	result := orm.Eloquent.Table("sys_menu").Create(&e)
	if result.Error != nil {
		err = result.Error
		return
	}
	err = InitPaths(e)
	if err != nil {
		return
	}
	id = e.MenuId
	return
}

func InitPaths(menu *Menu) (err error) {
	parentMenu := new(Menu)
	if int(menu.ParentId) != 0 {
		orm.Eloquent.Table("sys_menu").Where("menu_id = ?", menu.ParentId).First(parentMenu)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return
		}
		menu.Paths = parentMenu.Paths + "/" + utils.Int64ToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + utils.Int64ToString(menu.MenuId)
	}
	orm.Eloquent.Table("sys_menu").Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
	return
}

func (e *Menu) Update(id int64) (update Menu, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_menu").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_menu").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	err = InitPaths(e)
	if err != nil {
		return
	}
	return
}

func (e *Menu) Delete(id int64) (success bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_menu").Where("menu_id = ?", id).Update(mp).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
