package models

type RoleMenu struct {
	RoleId   int    `gorm:""`
	MenuId   int    `gorm:""`
	RoleName string `gorm:"size:128)"`
	ControlBy
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}
