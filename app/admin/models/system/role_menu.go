package system

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type RoleMenu struct {
	RoleId   int    `gorm:""`
	MenuId   int    `gorm:""`
	RoleName string `gorm:"size:128"`
	CreateBy string `gorm:"size:128"`
	UpdateBy string `gorm:"size:128"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

type MenuPath struct {
	Path string `json:"path"`
}

func (rm *RoleMenu) Get(tx *gorm.DB) ([]RoleMenu, error) {
	var r []RoleMenu
	table := tx.Table("sys_role_menu")
	if rm.RoleId != 0 {
		table = table.Where("role_id = ?", rm.RoleId)

	}
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) GetPermis(tx *gorm.DB) ([]string, error) {
	var r []SysMenu
	table := tx.Select("sys_menu.permission").Table("sys_menu").Joins("left join sys_role_menu on sys_menu.menu_id = sys_role_menu.menu_id")

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

func (rm *RoleMenu) GetIDS(tx *gorm.DB) ([]MenuPath, error) {
	var r []MenuPath
	table := tx.Select("sys_menu.path").Table("sys_role_menu")
	table = table.Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id")
	table = table.Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id")
	table = table.Where("sys_role.role_name = ? and sys_menu.type=1", rm.RoleName)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) DeleteRoleMenu(tx *gorm.DB, roleId int) error {
	if err := tx.Table("sys_role_dept").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		return err
	}
	if err := tx.Table("sys_role_menu").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		return err
	}
	var role SysRole
	if err := tx.Table("sys_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		return err
	}
	sql3 := "delete from sys_casbin_rule where v0= '" + role.RoleKey + "';"
	if err := tx.Exec(sql3).Error; err != nil {
		return err
	}
	return nil

}

// 该方法即将弃用
func (rm *RoleMenu) BatchDeleteRoleMenu(tx *gorm.DB, roleIds []int) error {
	if err := tx.Table("sys_role_menu").Where("role_id in (?)", roleIds).Delete(&rm).Error; err != nil {
		return err
	}
	var role []SysRole
	if err := tx.Table("sys_role").Where("role_id in (?)", roleIds).Find(&role).Error; err != nil {
		return err
	}
	sql := ""
	for i := 0; i < len(role); i++ {
		sql += "delete from sys_casbin_rule where v0= '" + role[i].RoleName + "';"
	}
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil

}

func (rm *RoleMenu) Insert(tx *gorm.DB, enforcer *casbin.SyncedEnforcer, roleId int, menuId []int) error {
	var err error
	var (
		role        SysRole
		menu        []SysMenu
		casbinRules []CasbinRule // casbinRule 待插入队列
	)
	// 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
	if err = tx.Table("sys_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		return err
	}
	if err = tx.Table("sys_menu").Where("menu_id in (?)", menuId).Find(&menu).Error; err != nil {
		return err
	}
	//ORM不支持批量插入所以需要拼接 sql 串
	sysRoleMenuSql := "INSERT INTO `sys_role_menu` (`role_id`,`menu_id`,`role_name`) VALUES "

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
			casbinRules = append(casbinRules,
				CasbinRule{
					V0: role.RoleKey,
					V1: m.Path,
					V2: m.Action,
				})
		}
	}
	// 执行批量插入sys_role_menu
	if err = tx.Exec(sysRoleMenuSql).Error; err != nil {
		return err
	}

	// 执行批量插入sys_casbin_rule
	if len(casbinRules) > 0 {
		if err = tx.Create(&casbinRules).Error; err != nil {
			return err
		}
	}
	return nil
}

func (rm *RoleMenu) Delete(tx *gorm.DB, RoleId string, MenuID string) (bool, error) {
	rm.RoleId, _ = pkg.StringToInt(RoleId)
	table := tx.Table("sys_role_menu").Where("role_id = ?", RoleId)
	if MenuID != "" {
		table = table.Where("menu_id = ?", MenuID)
	}
	if err := table.Delete(&rm).Error; err != nil {
		return false, err
	}
	return true, nil

}
