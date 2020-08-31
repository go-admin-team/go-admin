package model

import "gorm.io/gorm/schema"

type ActiveRecord interface {
	schema.Tabler
	SetCreateBy(createBy string)
	SetUpdateBy(updateBy string)
	Generate() ActiveRecord
	GetId() interface{}
}
