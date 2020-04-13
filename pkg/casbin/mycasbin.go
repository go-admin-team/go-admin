package mycasbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-sql-driver/mysql"
	"go-admin/config"
	"go-admin/pkg/utils"
)

var Em endpoint.Middleware

func Casbin() (*casbin.Enforcer, error) {
	conn := config.DatabaseConfig.Username + ":" + config.DatabaseConfig.Password + "@tcp(" + config.DatabaseConfig.Host + ":" + utils.IntToString(config.DatabaseConfig.Port) + ")/" + config.DatabaseConfig.Database
	if config.DatabaseConfig.Dbtype == "sqlite3" {
		conn = config.DatabaseConfig.Host
	}
	Apter, err := gormadapter.NewAdapter(config.DatabaseConfig.Dbtype, conn, true)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("config/rbac_model.conf", Apter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		fmt.Print("casbin rbac_model or policy init error, message: %v", err)
		return nil, err
	}
}
