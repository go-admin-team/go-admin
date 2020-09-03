package model

type ByModel struct {
	CreateBy uint `gorm:"index;"` // 创建者
	UpdateBy uint `gorm:"index;"` // 更新者
}
