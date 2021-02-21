package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
	"github.com/gin-gonic/gin"

	"go-admin/tools"
)

func LoadPolicy(c *gin.Context) (*casbin.SyncedEnforcer, error) {
	if err := Cfg.GetCasbinKey(c.Request.Host).LoadPolicy(); err == nil {
		return Cfg.GetCasbinKey(c.Request.Host), err
	} else {
		msgID := tools.GenerateMsgIDFromContext(c)
		log.LogPrintf("msgID[%s] casbin rbac_model or policy init error, message: %v \r\n", msgID, err.Error())
		return nil, err
	}
}
