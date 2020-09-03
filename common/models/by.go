package models

type ControlBy struct {
	CreateBy uint `gorm:"index;comment:'创建者'"`
	UpdateBy uint `gorm:"index;comment:'更新者'"`
}
