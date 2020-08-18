package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

var GinEngine *gin.Engine
var CasbinEnforcer *casbin.SyncedEnforcer
var Eloquent *gorm.DB

var GADMCron *cron.Cron

var (
	Source string
	Driver string
	DBName string
)

// go-admin Version Info
var Version string

func init() {
	Version = "1.1.5"
}

var (
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)
