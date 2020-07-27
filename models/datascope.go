package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-admin/tools"
	"go-admin/tools/config"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func (e *DataPermission) GetDataScope(tbname string, table *gorm.DB) (*gorm.DB, error) {

	if !config.ApplicationConfig.EnableDP {
		usageStr := `数据权限已经为您` + tools.Green(`关闭`) + `，如需开启请参考配置文件字段说明`
		fmt.Printf("%s\n", usageStr)
		return table, nil
	}
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
