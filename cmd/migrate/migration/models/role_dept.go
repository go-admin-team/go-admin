package models

//sys_role_dept
type SysRoleDept struct {
	RoleId int `gorm:"size:11;primaryKey"`
	DeptId int `gorm:"size:11;primaryKey"`
}

func (SysRoleDept) TableName() string {
	return "sys_role_dept"
}
