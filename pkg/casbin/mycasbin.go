package mycasbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

func Casbin() (*casbin.Enforcer, error) {
	Apter, err := gormAdapter.NewAdapterByDB(orm.Eloquent)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, Apter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		log.Printf("casbin rbac_model or policy init error, message: %v \r\n", err.Error())
		return nil, err
	}
}
