package dto

import (
	"github.com/go-admin-team/go-admin-core/v2/tools/search"
	"go-admin/common/global"
	"gorm.io/gorm"
)

type GeneralDelDto struct {
	Id  int   `uri:"id" json:"id" validate:"required"`
	Ids []int `json:"ids"`
}

func (g GeneralDelDto) GetIds() []int {
	ids := make([]int, 0)
	// Id 此前在 else 分支里被重复追加：仅传 Id 时会得到 [5 5]，
	// 同一条记录被执行两次删除
	if g.Id > 0 {
		ids = append(ids, g.Id)
	}
	for _, id := range g.Ids {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		//方式全部删除
		ids = append(ids, 0)
	}
	return ids
}

type GeneralGetDto struct {
	Id int `uri:"id" json:"id" validate:"required"`
}

func MakeCondition(q interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &search.GormCondition{
			GormPublic: search.GormPublic{},
			Join:       make([]*search.GormJoin, 0),
		}
		search.ResolveSearchQuery(global.Driver, q, condition)
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
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}
