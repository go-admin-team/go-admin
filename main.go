package main

import (
	"go-admin/cmd"
)

// @title go-admin API
// @version 1.0.1
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @description 添加qq群: 74520518 进入技术交流群 请备注，谢谢！
// @license.name MIT
// @license.url https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

//func main() {
//	configName := "settings"
//
//
//	config.InitConfig(configName)
//
//	gin.SetMode(gin.DebugMode)
//	log.Println(config.DatabaseConfig.Port)
//
//	err := gorm.AutoMigrate(orm.Eloquent)
//	if err != nil {
//		log.Fatalln("数据库初始化失败 err: %v", err)
//	}
//
//	if config.ApplicationConfig.IsInit {
//		if err := models.InitDb(); err != nil {
//			log.Fatal("数据库基础数据初始化失败！")
//		} else {
//			config.SetApplicationIsInit()
//		}
//	}
//
//	r := router.InitRouter()
//
//	defer orm.Eloquent.Close()
//
//	srv := &http.Server{
//		Addr:    config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
//		Handler: r,
//	}
//
//	go func() {
//		// 服务连接
//		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			log.Fatalf("listen: %s\n", err)
//		}
//	}()
//	log.Println("Server Run ", config.ApplicationConfig.Host+":"+config.ApplicationConfig.Port)
//	log.Println("Enter Control + C Shutdown Server")
//	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
//	quit := make(chan os.Signal)
//	signal.Notify(quit, os.Interrupt)
//	<-quit
//	log.Println("Shutdown Server ...")
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := srv.Shutdown(ctx); err != nil {
//		log.Fatal("Server Shutdown:", err)
//	}
//	log.Println("Server exiting")
//}

func main() {
	cmd.Execute()
}
