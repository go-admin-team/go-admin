package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/service"
	"go-admin/tools"
	"gorm.io/gorm"
)

type SysMenu struct {
	service.Service
}

// GetSysMenuPage 获取SysMenu列表
func (e *SysMenu) getSysMenuPage(c *dto.SysMenuSearch, list *[]models.SysMenu) error {
	var err error
	var data models.SysMenu
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysMenu 获取SysMenu对象
func (e *SysMenu) GetSysMenu(d *dto.SysMenuById, model *models.SysMenu) error {
	var err error
	var data models.SysMenu
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// InsertSysMenu 创建SysMenu对象
func (e *SysMenu) InsertSysMenu(model *models.SysMenu) error {
	var err error
	var data models.SysMenu
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

func (e *SysMenu) initPaths(menu *models.SysMenu) error {
	var err error
	var data models.SysMenu
	parentMenu := new(models.SysMenu)
	if int(menu.ParentId) != 0 {
		e.Orm.Model(&data).First(parentMenu, menu.ParentId)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return err
		}
		menu.Paths = parentMenu.Paths + "/" + tools.IntToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + tools.IntToString(menu.MenuId)
	}
	e.Orm.Model(&data).Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
	return err
}

// UpdateSysMenu 修改SysMenu对象
func (e *SysMenu) UpdateSysMenu(c *models.SysMenu) error {
	var err error
	var data models.SysMenu
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
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
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(d.GetId()).Delete(&data)
	if db.Error != nil {
		err = db.Error
		log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
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
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// SetSysMenuTree 设置菜单数据
func (e *SysMenu) SetSysMenuTree() (c *dto.SysMenuSearch, m []dto.MenuLabel, err error) {
	var list []models.SysMenu
	err = e.GetSysMenuList(c, &list)

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

// GetMenuList 菜单列表
func (e *SysMenu) GetSysMenuPage() (c *dto.SysMenuSearch, m []models.SysMenu, err error) {
	err = e.getSysMenuPage(c,&m)
	m = make([]models.SysMenu, 0)
	for i := 0; i < len(m); i++ {
		if m[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&m, m[i])

		m = append(m, menusInfo)
	}
	return
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