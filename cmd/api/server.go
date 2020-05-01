package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-admin/database"
	"go-admin/router"
	"go-admin/tools"
	config2 "go-admin/tools/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	config string
	port   string
	mode   string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "go-admin server config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func usage() {
	usageStr := `starting api server`
	fmt.Printf("%s\n", usageStr)
}

func setup() {

	//1. 读取配置
	config2.ConfigSetup(config)
	//2. 设置日志
	tools.InitLogger()
	//3. 初始化数据库链接
	database.Setup()
	//4. 设置gin mode
	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}
}

func run() error {
	r := router.InitRouter()

	defer database.Eloquent.Close()

	srv := &http.Server{
		Addr:    config2.ApplicationConfig.Host + ":" + config2.ApplicationConfig.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	content, _ := ioutil.ReadFile("./static/go-admin.txt")
	log.Println(string(content))
	log.Println("Server Run http://127.0.0.1:"+config2.ApplicationConfig.Port+"/")
	log.Println("Swagger URL http://127.0.0.1:"+config2.ApplicationConfig.Port+"/swagger/index.html")

	log.Println("Enter Control + C Shutdown Server")
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}
