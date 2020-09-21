package mycasbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"

	"go-admin/common/global"
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

func Setup() {
	Apter, err := gormAdapter.NewAdapterByDBUsePrefix(global.Eloquent, "sys_")
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewSyncedEnforcer(m, Apter)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
	global.CasbinEnforcer = e
}

func Casbin() *casbin.SyncedEnforcer {
	return global.CasbinEnforcer
}

func LoadPolicy() (*casbin.SyncedEnforcer, error) {
	if err := global.CasbinEnforcer.LoadPolicy(); err == nil {
		return global.CasbinEnforcer, err
	} else {
		log.Printf("casbin rbac_model or policy init error, message: %v \r\n", err.Error())
		return nil, err
	}
}
