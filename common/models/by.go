package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

type ModelTime struct {
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`

	// DeletedAt is milliseconds since the epoch, zero while the row is live,
	// and never null.
	//
	// A nullable marker cannot take part in a unique index. Two live rows are
	// (name, NULL) and (name, NULL), and NULL is not equal to NULL, so the
	// index permits both — it looks like a constraint and enforces nothing.
	// With zero for live rows the pair collides, while two deletions of the
	// same name differ by their timestamps and both remain.
	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"softDelete:milli;index;comment:删除时间"`
}
