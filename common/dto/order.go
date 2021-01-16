package dto

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func OrderDest(sort string, bl bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: bl})
	}
}
