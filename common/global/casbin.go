package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go-admin/common/apis"
	"go-admin/tools/app"
)

func LoadPolicy(c *gin.Context) (*casbin.SyncedEnforcer, error) {
	log := apis.GetRequestLogger(c)
	if err := app.Runtime.GetCasbinKey(c.Request.Host).LoadPolicy(); err == nil {
		return app.Runtime.GetCasbinKey(c.Request.Host), err
	} else {
		log.Errorf("casbin rbac_model or policy init error, %s ", err.Error())
		return nil, err
	}
}
