package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wenjianzhang/go-admin/global"
	mycasbin "github.com/wenjianzhang/go-admin/pkg/casbin"
	"github.com/wenjianzhang/go-admin/pkg/logger"
	"github.com/wenjianzhang/go-admin/router"
	"log"
)

func main() {
	var err error
	global.Eloquent, err = gorm.Open("mysql", "root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	mycasbin.Setup()
	logger.Setup()
	global.GinEngine = gin.Default()
	router.InitRouter()
	log.Fatal(global.GinEngine.Run(":8000"))
}
