package models

import "time"

type ControlBy struct {
	CreateBy int `gorm:"index;comment:'创建者'"`
	UpdateBy int `gorm:"index;comment:'更新者'"`
}

func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}


type Model struct {
	ID int `gorm:"primaryKey;autoIncrement;comment:'主键编码'"`
}

type ModelTime struct {
	CreatedAt time.Time  `gorm:"comment:'创建时间'"`
	UpdatedAt time.Time  `gorm:"comment:'最后更新时间'"`
	DeletedAt *time.Time `gorm:"index;comment:'删除时间'"`
}
