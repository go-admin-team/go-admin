package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
)

func LoadPolicy() (*casbin.SyncedEnforcer, error) {
	if err := CasbinEnforcer.LoadPolicy(); err == nil {
		return CasbinEnforcer, err
	} else {
		log.LogPrintf("casbin rbac_model or policy init error, message: %v \r\n", err.Error())
		return nil, err
	}
}
