package models

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
