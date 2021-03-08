package service

import (
	log "github.com/go-admin-team/go-admin-core/logger"
	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysRoleMenu struct {
	service.Service
}

func (e *SysRoleMenu) ReloadRule(tx *gorm.DB, roleId int, menuId []int) (err error) {
	var role system.SysRole

	msgID := e.MsgID

	menu := make([]models.Menu, 0)
	roleMenu := make([]system.RoleMenu, len(menuId))
	casbinRule := make([]system.CasbinRule, 0)
	//先删除所有的
	err = e.DeleteRoleMenu(tx, roleId)
	if err != nil {
		return
	}

	// 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
	err = tx.Where("role_id = ?", roleId).First(&role).Error
	if err != nil {
		log.Errorf("msgID[%s] get role error, %s", msgID, err.Error())
		return
	}
	err = tx.Where("menu_id in (?)", menuId).
		//Select("path, action, menu_id, menu_type").
		Find(&menu).Error
	if err != nil {
		log.Errorf("msgID[%s] get menu error, %s", msgID, err.Error())
		return
	}
	for i := range menu {
		roleMenu[i] = system.RoleMenu{
			RoleId:   role.RoleId,
			MenuId:   menu[i].MenuId,
			RoleName: role.RoleKey,
		}
		if menu[i].MenuType == "A" {
			casbinRule = append(casbinRule, system.CasbinRule{
				PType: "p",
				V0:    role.RoleKey,
				V1:    menu[i].Path,
				V2:    menu[i].Action,
			})
		}
	}
	err = tx.Create(&roleMenu).Error
	if err != nil {
		log.Errorf("msgID[%s] batch create role's menu error, %s", msgID, err.Error())
		return
	}
	if len(casbinRule) > 0 {
		err = tx.Create(&casbinRule).Error
		if err != nil {
			log.Errorf("msgID[%s] batch create casbin rule error, %s", msgID, err.Error())
			return
		}
	}

	return
}

func (e *SysRoleMenu) DeleteRoleMenu(tx *gorm.DB, roleId int) (err error) {
	msgID := e.MsgID
	err = tx.Where("role_id = ?", roleId).
		Delete(&system.SysRoleDept{}).Error
	if err != nil {
		log.Errorf("msgID[%s] delete role's dept error, %s", msgID, err.Error())
		return
	}
	err = tx.Where("role_id = ?", roleId).
		Delete(&system.RoleMenu{}).Error
	if err != nil {
		log.Errorf("msgID[%s] delete role's menu error, %s", msgID, err.Error())
		return
	}
	var role system.SysRole
	err = tx.Where("role_id = ?", roleId).
		First(&role).Error
	if err != nil {
		log.Errorf("msgID[%s] get role error, %s", msgID, err.Error())
		return
	}
	err = tx.Where("v0 = ?", role.RoleKey).
		Delete(&system.CasbinRule{}).Error
	if err != nil {
		log.Errorf("msgID[%s] delete casbin rule error, %s", msgID, err.Error())
		return
	}
	return
}

func (e *SysRoleMenu) GetIDS(tx *gorm.DB, roleName string) ([]system.MenuPath, error) {
	var r []system.MenuPath
	table := tx.Select("sys_menu.path").Table("sys_role_menu")
	table = table.Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id")
	table = table.Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id")
	table = table.Where("sys_role.role_name = ? and sys_menu.type=1", roleName)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}
