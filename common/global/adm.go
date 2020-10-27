package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"go-admin/common/config"
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
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)
