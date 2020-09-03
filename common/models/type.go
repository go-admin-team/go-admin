package models

import "gorm.io/gorm/schema"

type ActiveRecord interface {
	schema.Tabler
	SetCreateBy(createBy uint)
	SetUpdateBy(updateBy uint)
	Generate() ActiveRecord
	GetId() interface{}
}
