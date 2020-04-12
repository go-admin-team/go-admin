package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/config"
	orm "go-admin/database"
	"go-admin/models"
	"go-admin/router"
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

	server := &http.Server{
		Addr:              config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler:           r,
	}

	go func() {
		err := server.ListenAndServe(); if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
