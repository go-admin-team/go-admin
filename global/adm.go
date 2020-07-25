package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/jinzhu/gorm"
)

var GinEngine *gin.Engine
var CasbinEnforcer *casbin.Enforcer
var Eloquent *gorm.DB

var Source string
var Driver string
var DBName string

// go-admin Version Info
var Version string

func init() {
	Version = "1.0.10"
}

var Logger *glog.Logger
var DBLogger *glog.Logger
var AccessLogger *glog.Logger
