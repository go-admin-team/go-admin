package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"github.com/spf13/cobra"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/router"
	"go-admin/app/jobs"
	"go-admin/common/database"
	"go-admin/common/global"
	ext "go-admin/config"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
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

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {
	// 注入配置扩展项
	config.ExtendConfig = &ext.ExtConfig
	//1. 读取配置
	config.Setup(file.NewSource, file.WithPath(configYml))
	go config.Watch()
	//2. 设置日志
	sdk.Runtime.SetLogger(
		logger.SetupLogger(
			config.LoggerConfig.Type,
			config.LoggerConfig.Path,
			config.LoggerConfig.Level,
			config.LoggerConfig.Stdout))
	//3. 初始化数据库链接
	database.Setup()
	//4. 设置缓存
	cacheAdapter, err := config.CacheConfig.Setup()
	if err != nil {
		log.Fatalf("cache setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetCacheAdapter(cacheAdapter)
	//5. 设置验证码store
	captcha.SetStore(captcha.NewCacheStore(cacheAdapter, 600))

	//6. 设置队列
	if !config.QueueConfig.Empty() {
		queueAdapter, err := config.QueueConfig.Setup()
		if err != nil {
			log.Fatalf("queue setup error, %s\n", err.Error())
		}
		sdk.Runtime.SetQueueAdapter(queueAdapter)
	}

	//7. 设置分布式锁
	if !config.LockerConfig.Empty() {
		lockerAdapter, err := config.LockerConfig.Setup()
		if err != nil {
			log.Fatalf("locker setup error, %s\n", err.Error())
		}
		sdk.Runtime.SetLockerAdapter(lockerAdapter)
	}

	usageStr := `starting api server...`
	log.Println(usageStr)
	//注册监听函数
	queue := sdk.Runtime.GetMemoryQueue("")
	queue.Register(global.LoginLog, system.SaveLoginLog)
	queue.Register(global.OperateLog, system.SaveOperaLog)
	go queue.Run()
}

func run() error {
	defer config.Stop()

	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := sdk.Runtime.GetEngine()
	if engine == nil {
		engine = gin.New()
	}

	if config.ApplicationConfig.Mode == "dev" {
		//监控
		AppRouters = append(AppRouters, router.Monitor)
	}

	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler: sdk.Runtime.GetEngine(),
	}
	go func() {
		jobs.InitJob()
		jobs.Setup(sdk.Runtime.GetDb())

	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	fmt.Println(pkg.Red(string(global.LogoContent)))
	tip()
	fmt.Println(pkg.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/ \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%d/ \r\n", pkg.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Println(pkg.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/swagger/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%d/swagger/index.html \r\n", pkg.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", pkg.GetCurrentTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", pkg.GetCurrentTimeStr())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	return nil
}

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`go-admin `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}
