package models

import (
	"fmt"
	orm "go-admin/database"
)

//sys_role_dept
type SysRoleDept struct {
	RoleId int64
	DeptId int64
}

func (rm *SysRoleDept) Insert(roleId int64, deptIds []int64) (bool, error) {
	//ORM不支持批量插入所以需要拼接 sql 串
	sql := "INSERT INTO `sys_role_dept` (`role_id`,`dept_id`) VALUES "

	for i := 0; i < len(deptIds); i++ {
		if len(deptIds)-1 == i {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,%d);", roleId, deptIds[i])
		} else {
			sql += fmt.Sprintf("(%d,%d),", roleId, deptIds[i])
		}
	}
	orm.Eloquent.Exec(sql)

	return true, nil
}

func (rm *SysRoleDept) DeleteRoleDept(roleId int64) (bool, error) {
	if err := orm.Eloquent.Table("sys_role_dept").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		return false, err
	}
	var role SysRole
	if err := orm.Eloquent.Table("sys_role").Where("id = ?", roleId).First(&role).Error; err != nil {
		return false, err
	}

	return true, nil

}
