package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"go-admin/common/config"
	"go-admin/pkg/logger"
)

const (
	// go-admin Version Info
	Version = "1.2.2"
)

var Cfg config.Conf = config.DefaultConfig()

var GinEngine *gin.Engine
var CasbinEnforcer *casbin.SyncedEnforcer
var Eloquent *gorm.DB

var GADMCron *cron.Cron

var (
	Source string
	Driver string
	DBName string
)

var (
	Logger        = &logger.Logger{}
	JobLogger     = &logger.Logger{}
	RequestLogger = &logger.Logger{}
)
