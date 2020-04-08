package models

import (
	"github.com/jinzhu/gorm"
	"go-admin/utils"
)

type DataPermission struct {
	DataScope string
	UserId    int64
	DeptId    int64
	RoleId    int64
}

func (e *DataPermission) GetDataScope(tbname string, table *gorm.DB) *gorm.DB {
	SysUser := new(SysUser)
	SysRole := new(SysRole)
	SysUser.Id = e.UserId
	user, _ := SysUser.Get()
	SysRole.Id = user.RoleId
	role, _ := SysRole.Get()

	if role.DataScope == "2" {
		table = table.Where(tbname+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		table = table.Where(tbname+".create_by in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		table = table.Where(tbname+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+utils.Int64ToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		table = table.Where(tbname+".create_by = ?", e.UserId)
	}

	return table
}
