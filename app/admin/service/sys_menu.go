package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysMenu struct {
	service.Service
}

// GetSysMenuPage 获取SysMenu列表
func (e *SysMenu) GetSysMenuPage(c *dto.SysMenuSearch) (*[]models.SysMenu, error) {
	var m = make([]models.SysMenu, 0)
	var err error
	var menu = make([]models.SysMenu, 0)
	err = e.getSysMenuPage(c, &menu)
	for i := 0; i < len(menu); i++ {
		if menu[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menu, menu[i])

		m = append(m, menusInfo)
	}
	return &m, err
}

// getSysMenuPage 菜单列表
func (e *SysMenu) getSysMenuPage(c *dto.SysMenuSearch, list *[]models.SysMenu) error {
	var err error
	var data models.SysMenu

	err = e.Orm.Model(&data).Preload("SysApi").
		Scopes(
			cDto.OrderDest("sort", false),
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetSysMenu 获取SysMenu对象
func (e *SysMenu) GetSysMenu(d *dto.SysMenuById, model *models.SysMenu) error {
	var err error
	var data models.SysMenu

	db := e.Orm.Model(&data).Preload("SysApi").
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	for _, v := range model.SysApi {
		model.Apis = append(model.Apis, v.Id)
	}
	return nil
}

// InsertSysMenu 创建SysMenu对象
func (e *SysMenu) InsertSysMenu(model *models.SysMenu) error {
	var err error
	var data models.SysMenu

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *SysMenu) initPaths(menu *models.SysMenu) error {
	var err error
	var data models.SysMenu
	parentMenu := new(models.SysMenu)
	if menu.ParentId != 0 {
		e.Orm.Model(&data).First(parentMenu, menu.ParentId)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return err
		}
		menu.Paths = parentMenu.Paths + "/" + pkg.IntToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + pkg.IntToString(menu.MenuId)
	}
	e.Orm.Model(&data).Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
	return err
}

// UpdateSysMenu 修改SysMenu对象
func (e *SysMenu) UpdateSysMenu(c *models.SysMenu) error {
	var err error

	db := e.Orm.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(c)
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}          

// RemoveSysMenu 删除SysMenu
func (e *SysMenu) RemoveSysMenu(d *dto.SysMenuById) error {
	var err error
	var data models.SysMenu

	db := e.Orm.Model(&data).
		Where(d.Ids).Delete(&data)
	if db.Error != nil {
		err = db.Error
		e.Log.Errorf("Delete error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}

// GetSysMenuList 获取菜单数据
func (e *SysMenu) GetSysMenuList(c *dto.SysMenuSearch, list *[]models.SysMenu) error {
	var err error
	var data models.SysMenu

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// SetSysMenuTree 设置菜单数据
func (e *SysMenu) SetSysMenuLabel() (m []dto.MenuLabel, err error) {
	var list []models.SysMenu
	err = e.GetSysMenuList(&dto.SysMenuSearch{}, &list)

	m = make([]dto.MenuLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.MenuLabel{}
		e.Id = list[i].MenuId
		e.Label = list[i].Title
		deptsInfo := menuLabelCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// 左侧菜单
func (e *SysMenu) GetSysMenuByRoleName(roleName string) (Menus []models.SysMenu, err error) {
	var table *gorm.DB
	var data models.SysMenu

	if roleName == "admin" {
		table = e.Orm.Model(&data).Select("sys_menu.*")
		table = table.Where(" menu_type in ('M','C')")
	} else {
		table = e.Orm.Model(&data).Select("sys_menu.*").Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id")
		table = table.Where("sys_role_menu.role_name=? and menu_type in ('M','C')", roleName)
	}

	err = table.Order("sort").Find(&Menus).Error

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return
}

// menuLabelCall 递归构造组织数据
func menuLabelCall(elist *[]models.SysMenu, dept dto.MenuLabel) dto.MenuLabel {
	list := *elist

	min := make([]dto.MenuLabel, 0)
	for j := 0; j < len(list); j++ {

		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.MenuLabel{}
		mi.Id = list[j].MenuId
		mi.Label = list[j].Title
		mi.Children = []dto.MenuLabel{}
		if list[j].MenuType != "F" {
			ms := menuLabelCall(elist, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	if len(min) > 0 {
		dept.Children = min
	} else {
		dept.Children = nil
	}
	return dept
}

func menuCall(menulist *[]models.SysMenu, menu models.SysMenu) models.SysMenu {
	list := *menulist

	min := make([]models.SysMenu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuId != list[j].ParentId {
			continue
		}
		mi := models.SysMenu{}
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
		mi.Visible = list[j].Visible
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []models.SysMenu{}

		if mi.MenuType != "F" {
			ms := menuCall(menulist, mi)
			min = append(min, ms)

		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}

func (e *SysMenu) SetMenuRole(roleName string) (m []models.SysMenu, err error) {

	menus, err := e.getByRoleName(roleName)

	m = make([]models.SysMenu, 0)
	for i := 0; i < len(menus); i++ {
		if menus[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menus, menus[i])

		m = append(m, menusInfo)
	}
	return
}

func (e *SysMenu) getByRoleName(roleName string) (Menus []models.SysMenu, err error) {
	var data models.SysMenu

	var table *gorm.DB
	if roleName == "admin" {
		table = e.Orm.Model(&data).Select("sys_menu.*")
		table = table.Where(" menu_type in ('M','C')")
	} else {
		table = e.Orm.Model(&data).Select("sys_menu.*").Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id")
		table = table.Where("sys_role_menu.role_name=? and menu_type in ('M','C')", roleName)
	}
	err = table.Scopes(
		cDto.OrderDest("sort", false),
	).Find(&Menus).Error

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return
}