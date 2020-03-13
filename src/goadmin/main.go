package main

import (
	"github.com/gin-gonic/gin"
	"goadmin/config"
	orm "goadmin/database"
	"goadmin/models"
	"goadmin/router"
	"log"
)

// @title goadmin API
// @version 0.0.1
// @description Swagger.

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
