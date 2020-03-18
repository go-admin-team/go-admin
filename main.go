package main

import (
	"github.com/gin-gonic/gin"
	"go-admin/config"
	orm "go-admin/database"
	"go-admin/models"
	"go-admin/router"
	"log"
)

// @title go-admin API
// @version 0.0.1
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @description 添加qq群: 74520518 进入技术交流群 请备注，谢谢！
// @license.name MIT
// @license.url https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {

	gin.SetMode(gin.DebugMode)

	log.Println(config.DatabaseConfig.Port)
	if config.ApplicationConfig.IsInit {
		if err := models.InitDb(); err != nil {
			log.Fatal("数据库初始化失败！")
		} else {
			config.SetApplicationIsInit()
		}
	}
	r := router.InitRouter()

	defer orm.Eloquent.Close()
	if err := r.Run(config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port); err != nil {
		log.Fatal(err)
	}

}
