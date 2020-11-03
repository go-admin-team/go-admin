package api

import (
	"context"
	"fmt"
	"go-admin/tools/trace"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-admin/app/admin/router"
	"go-admin/app/jobs"
	"go-admin/common/database"
	"go-admin/common/global"
	"go-admin/common/log"
	mycasbin "go-admin/pkg/casbin"
	"go-admin/pkg/logger"
	"go-admin/tools"
	"go-admin/tools/config"
)

var (
	configYml  string
	port       string
	mode       string
	traceStart bool
	StartCmd   = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "go-admin server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var AppRouters = make([]func(), 0)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
	StartCmd.PersistentFlags().BoolVarP(&traceStart, "traceStart", "t", false, "start traceStart app dash")

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {

	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	global.Logger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "bus")
	global.JobLogger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "job")
	global.RequestLogger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "request")
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)
	//4. 接口访问控制加载
	global.CasbinEnforcer = mycasbin.Setup(global.Eloquent, "sys_")

	usageStr := `starting api server`
	log.Info(usageStr)

}

func run() error {
	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := global.Cfg.GetEngine()
	if engine == nil {
		engine = gin.New()
	}

	if mode == "dev" {
		//监控
		AppRouters = append(AppRouters, router.Monitor)
	}

	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:    config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler: global.Cfg.GetEngine(),
	}
	go func() {
		jobs.InitJob()
		jobs.Setup()

	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if traceStart {
		//链路追踪, fixme 页面显示需要自备梯子
		trace.Start()
		defer trace.Stop(ctx)
	}

	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("listen: ", err)
			}
		}
	}()
	fmt.Println(tools.Red(string(global.LogoContent)))
	tip()
	fmt.Println(tools.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/ \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/ \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Println(tools.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/swagger/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/swagger/index.html \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", tools.GetCurrentTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrentTimeStr())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")

	return nil
}

func tip() {
	usageStr := `欢迎使用 ` + tools.Green(`go-admin `+global.Version) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}
