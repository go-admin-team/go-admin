package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-admin/tools"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func (e *DataPermission) GetDataScope(tbname string, table *gorm.DB) (*gorm.DB, error) {
	SysUser := new(SysUser)
	SysRole := new(SysRole)
	SysUser.UserId = e.UserId
	user, err := SysUser.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	SysRole.RoleId = user.RoleId
	role, err := SysRole.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
		table = table.Where(tbname+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		table = table.Where(tbname+".create_by in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		table = table.Where(tbname+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+tools.IntToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		table = table.Where(tbname+".create_by = ?", e.UserId)
	}

	return table, nil
}
