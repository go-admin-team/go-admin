package service

import (
	"errors"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	orm "go-admin/common/global"
	"go-admin/common/log"
	"go-admin/common/service"
)

type SysRole struct {
	service.Service
}

// GetSysRolePage 获取SysRole列表
func (e *SysRole) GetSysRolePage(c *dto.SysRoleSearch, list *[]system.SysRole, count *int64) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysRole 获取SysRole对象
func (e *SysRole) GetSysRole(d *dto.SysRoleById, model *system.SysRole) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	data.MenuIds, err = e.GetRoleMenuId(data.RoleId)
	if err != nil {
		log.Errorf("msgID[%s] get menuIds error, %s", msgID, err.Error())
		return err
	}
	return nil
}

// InsertSysRole 创建SysRole对象
func (e *SysRole) InsertSysRole(c *system.SysRole) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Model(&data).
		Create(c).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if len(c.MenuIds) > 0 {
		s := SysRoleMenu{}
		s.Orm = e.Orm
		s.MsgID = msgID
		err = s.ReloadRule(tx, c.RoleId, c.MenuIds)
		if err != nil {
			log.Errorf("msgID[%s] reload casbin rule error, %", msgID, err.Error())
			return err
		}
	}
	return nil
}

// UpdateSysRole 修改SysRole对象
func (e *SysRole) UpdateSysRole(c *system.SysRole) error {
	var err error
	msgID := e.MsgID

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	db := tx.Model(&c).
		Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	var t system.RoleMenu
	err = t.DeleteRoleMenu(tx, c.RoleId)
	if err != nil {
		log.Errorf("msgID[%s] delete role menu error, %", msgID, err.Error())
		return err
	}
	if len(c.MenuIds) > 0 {
		s := SysRoleMenu{}
		s.Orm = e.Orm
		s.MsgID = msgID
		err = s.ReloadRule(tx, c.RoleId, c.MenuIds)
		if err != nil {
			log.Errorf("msgID[%s] reload casbin rule error, %", msgID, err.Error())
			return err
		}
	}
	return nil
}

// RemoveSysRole 删除SysRole
func (e *SysRole) RemoveSysRole(d *dto.SysRoleById) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	db := tx.Model(&data).Delete(&data, d.GetId())
	if db.Error != nil {
		err = db.Error
		log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	s := SysRoleMenu{}
	s.Orm = db
	s.MsgID = msgID
	err = s.DeleteRoleMenu(tx, d.Id)
	if err != nil {
		log.Errorf("msgID[%s] insert role menu error, %", msgID, err.Error())
		return err
	}
	return nil
}

// 获取角色对应的菜单ids
func (e *SysRole) GetRoleMenuId(roleId int) ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]models.MenuIdList, 0)
	if err := orm.Eloquent.Table("sys_role_menu").
		Select("sys_role_menu.menu_id").
		Where("role_id = ? ", roleId).
		Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu "+
			"LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id =?  and parent_id is not null)", roleId).
		Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (e *SysRole) UpdateDataScope(c *system.SysRole) (err error) {
	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Model(&system.SysRole{}).Where("role_id = ?", c.RoleId).Select("data_scope, update_by").Updates(c).Error
	if err != nil {
		return err
	}

	err = tx.Where("role_id = ?", c.RoleId).Delete(&system.SysRoleDept{}).Error
	if err != nil {
		return err
	}

	if c.DataScope == "2" {
		deptRoles := make([]system.SysRoleDept, len(c.DeptIds))
		for i := range c.DeptIds {
			deptRoles[i] = system.SysRoleDept{
				RoleId: c.RoleId,
				DeptId: c.DeptIds[i],
			}
		}
		err = tx.Create(&deptRoles).Error
		if err != nil {
			return err
		}
	}
	return err
}
