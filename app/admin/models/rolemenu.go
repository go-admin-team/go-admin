package models

import (
	"fmt"

	orm "go-admin/common/global"
	"go-admin/tools"
)

type RoleMenu struct {
	RoleId   int    `gorm:""`
	MenuId   int    `gorm:""`
	RoleName string `gorm:"size:128)"`
	CreateBy string `gorm:"size:128)"`
	UpdateBy string `gorm:"size:128)"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

type MenuPath struct {
	Path string `json:"path"`
}

func (rm *RoleMenu) Get() ([]RoleMenu, error) {
	var r []RoleMenu
	table := orm.Eloquent.Table("sys_role_menu")
	if rm.RoleId != 0 {
		table = table.Where("role_id = ?", rm.RoleId)

	}
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) GetPermis() ([]string, error) {
	var r []Menu
	table := orm.Eloquent.Select("sys_menu.permission").Table("sys_menu").Joins("left join sys_role_menu on sys_menu.menu_id = sys_role_menu.menu_id")

	table = table.Where("role_id = ?", rm.RoleId)

	table = table.Where("sys_menu.menu_type in('F','C')")
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	var list []string
	for i := 0; i < len(r); i++ {
		list = append(list, r[i].Permission)
	}
	return list, nil
}

func (rm *RoleMenu) GetIDS() ([]MenuPath, error) {
	var r []MenuPath
	table := orm.Eloquent.Select("sys_menu.path").Table("sys_role_menu")
	table = table.Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id")
	table = table.Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id")
	table = table.Where("sys_role.role_name = ? and sys_menu.type=1", rm.RoleName)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) DeleteRoleMenu(roleId int) (bool, error) {
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := tx.Table("sys_role_dept").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Table("sys_role_menu").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	var role SysRole
	if err := tx.Table("sys_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	sql3 := "delete from sys_casbin_rule where v0= '" + role.RoleKey + "';"
	if err := tx.Exec(sql3).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil

}

// 该方法即将弃用
func (rm *RoleMenu) BatchDeleteRoleMenu(roleIds []int) (bool, error) {
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := tx.Table("sys_role_menu").Where("role_id in (?)", roleIds).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	var role []SysRole
	if err := tx.Table("sys_role").Where("role_id in (?)", roleIds).Find(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	sql := ""
	for i := 0; i < len(role); i++ {
		sql += "delete from sys_casbin_rule where v0= '" + role[i].RoleName + "';"
	}
	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil

}

func (rm *RoleMenu) Insert(roleId int, menuId []int) (bool, error) {
	var (
		role            SysRole
		menu            []Menu
		casbinRuleQueue []CasbinRule // casbinRule 待插入队列
	)

	// 开始事务
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	// 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
	if err := tx.Table("sys_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Table("sys_menu").Where("menu_id in (?)", menuId).Find(&menu).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	//ORM不支持批量插入所以需要拼接 sql 串
	sysRoleMenuSql := "INSERT INTO `sys_role_menu` (`role_id`,`menu_id`,`role_name`) VALUES "
	casbinRuleSql := "INSERT INTO `sys_casbin_rule`  (`p_type`,`v0`,`v1`,`v2`) VALUES "

	for i, m := range menu {
		// 拼装'role_menu'表批量插入SQL语句
		sysRoleMenuSql += fmt.Sprintf("(%d,%d,'%s')", role.RoleId, m.MenuId, role.RoleKey)
		if i == len(menu)-1 {
			sysRoleMenuSql += ";" //最后一条数据 以分号结尾
		} else {
			sysRoleMenuSql += ","
		}
		if m.MenuType == "A" {
			// 加入队列
			casbinRuleQueue = append(casbinRuleQueue,
				CasbinRule{
					V0: role.RoleKey,
					V1: m.Path,
					V2: m.Action,
				})
		}
	}
	// 执行批量插入sys_role_menu
	if err := tx.Exec(sysRoleMenuSql).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 拼装'sys_casbin_rule'批量插入SQL语句
	// TODO: casbinRuleQueue队列不为空时才会拼装，否则直接忽略不执行'for'循环
	for i, v := range casbinRuleQueue {
		casbinRuleSql += fmt.Sprintf("('p','%s','%s','%s')", v.V0, v.V1, v.V2)
		if i == len(casbinRuleQueue)-1 {
			casbinRuleSql += ";"
		} else {
			casbinRuleSql += ","
		}
	}
	// 执行批量插入sys_casbin_rule
	if len(casbinRuleQueue) > 0 {
		if err := tx.Exec(casbinRuleSql).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil
}

func (rm *RoleMenu) Delete(RoleId string, MenuID string) (bool, error) {
	rm.RoleId, _ = tools.StringToInt(RoleId)
	table := orm.Eloquent.Table("sys_role_menu").Where("role_id = ?", RoleId)
	if MenuID != "" {
		table = table.Where("menu_id = ?", MenuID)
	}
	if err := table.Delete(&rm).Error; err != nil {
		return false, err
	}
	return true, nil

}
