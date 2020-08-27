package search

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	FromQueryTag = "search"
	Mysql        = "mysql"    //mysql数据库标识
	Postgres     = "postgres" //postgres数据库标识
)

/**
 * 	exact / iexact 等于
 * 	contains / icontains 包含
 *	gt / gte 大于 / 大于等于
 *	lt / lte 小于 / 小于等于
 *	startswith / istartswith 以…起始
 *	endswith / iendswith 以…结束
 *	in
 *	isnull
 *  order 排序		e.g. order[key]=desc     order[key]=asc
 */
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	var tag string
	var ok bool
	var t resolveSearchTag
	for i := 0; i < qType.NumField(); i++ {
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(FromQueryTag)
		if !ok {
			//递归调用
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), condition)
			continue
		}
		switch tag {
		case "-":
			continue
		}
		t = resolveTagValue(tag)
		if qValue.Field(i).IsZero() {
			continue
		}
		//解析
		switch t.Type {
		case "left":
			//左关联
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"left join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case "exact", "iexact":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` = ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "contains", "icontains":
			//注意:mysql不支持ilike
			if driver == Postgres && t.Type == "icontains" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			}
		case "gt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` > ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "gte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` >= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` < ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` <= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "startswith", "istartswith":
			if driver == Postgres && t.Type == "istartswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			}
		case "endswith", "iendswith":
			if driver == Postgres && t.Type == "iendswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			}
		case "in":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` in (?)", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "isnull":
			if !(qValue.Field(i).IsZero() && qValue.Field(i).IsNil()) {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` isnull", t.Table, t.Column), make([]interface{}, 0))
			}
		case "order":
			switch strings.ToLower(qValue.Field(i).String()) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("`%s`.`%s` %s", t.Table, t.Column, qValue.Field(i).String()))
			}
		}
	}
	//return condition, order
}
