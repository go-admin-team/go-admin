package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go-admin/pkg/logger"
)

func LoadPolicy(c *gin.Context) (*casbin.SyncedEnforcer, error) {
	log := logger.GetRequestLogger(c)
	if err := Runtime.GetCasbinKey(c.Request.Host).LoadPolicy(); err == nil {
		return Runtime.GetCasbinKey(c.Request.Host), err
	} else {
		log.Errorf("casbin rbac_model or policy init error, %s ", err.Error())
		return nil, err
	}
}
