package models

import (
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

// ModelTime is frozen at the schema shape these tables had before
// 1786700003000 converted deleted_at to a NOT NULL millisecond marker. That is
// correct for the migrations ordered before the conversion, and wrong for any
// added after it: writes put NULL into a NOT NULL column, and reads are scoped
// "WHERE deleted_at IS NULL" and match nothing.
//
// Migrations after that version seed through the runtime models in app/.
// TestPostConversionMigrationsAvoidFrozenSeedModels enforces this.
type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
