package tools

import (
	"go-admin/tools/config"

	"github.com/matchstalk/utils/search"
	"gorm.io/gorm"
)

type GeneralDelDto struct {
	Id string `uri:"id" json:"id" validate:"required"`
}
type GeneralGetDto struct {
	Id string `uri:"id" json:"id" validate:"required"`
}

func SetQuery(db *gorm.DB, q interface{}) *gorm.DB {
	condition := &search.GormCondition{
		GormPublic: search.GormPublic{},
		Join:       make([]*search.GormJoin, 0),
	}
	search.ResolveSearchQuery(config.DatabaseConfig.Driver, q, condition)
	for _, join := range condition.Join {
		if join == nil {
			continue
		}
		db = db.Joins(join.JoinOn)
		for k, v := range join.Where {
			db = db.Where(k, v...)
		}
		for k, v := range join.Or {
			db = db.Or(k, v...)
		}
		for _, o := range join.Order {
			db = db.Order(o)
		}
	}
	for k, v := range condition.Where {
		db = db.Where(k, v...)
	}
	for k, v := range condition.Or {
		db = db.Or(k, v...)
	}
	for _, o := range condition.Order {
		db = db.Order(o)
	}
	return db
}
