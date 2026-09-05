package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/config/source/file"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go-admin/app/admin/models"
	"go-admin/app/admin/router"
	"go-admin/app/jobs"
	"go-admin/common/database"
	"go-admin/common/global"
	common "go-admin/common/middleware"
	"go-admin/common/middleware/handler"
	"go-admin/common/storage"
	ext "go-admin/config"
)

var (
	configYml string
	apiCheck  bool
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
	StartCmd.PersistentFlags().BoolVarP(&apiCheck, "api", "a", false, "Start server with check api data")

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {
	// 注入配置扩展项
	config.ExtendConfig = &ext.ExtConfig
	//1. 读取配置
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		storage.Setup,
	)
	//注册监听函数
	queue := sdk.Runtime.GetQueuePrefix("")
	queue.Register(global.LoginLog, models.SaveLoginLog)
	queue.Register(global.OperateLog, models.SaveOperaLog)
	queue.Register(global.ApiCheck, models.SaveSysApi)
	go queue.Run()

	usageStr := `starting api server...`
	log.Info(usageStr)
}

func run() error {
	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	initRouter()

	runStartupHooks()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler:      sdk.Runtime.GetEngine(),
		ReadTimeout:  time.Duration(config.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ApplicationConfig.WriterTimeout) * time.Second,
	}

	go func() {
		jobs.InitJob()
		jobs.Setup(sdk.Runtime.GetAllDb())

	}()

	if apiCheck {
		var routers = sdk.Runtime.GetRouter()
		q := sdk.Runtime.GetQueuePrefix("")
		mp := make(map[string]interface{})
		mp["List"] = routers
		message, err := sdk.Runtime.GetStreamMessage("", global.ApiCheck, mp)
		if err != nil {
			log.Infof("GetStreamMessage error, %s \n", err.Error())
			//日志报错错误，不中断请求
		} else {
			err = q.Append(message)
			if err != nil {
				log.Infof("Append message error, %s \n", err.Error())
			}
		}
	}

	// Armed before the server starts serving, and well before the readiness
	// banner: a signal arriving between "the process is up" and "the process
	// is listening for signals" reaches the default handler and kills it
	// without any of the shutdown below. That window is the whole reason
	// arming is separate from waiting.
	quit, disarmStopSignals := armStopSignals()

	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		}
	}()
	fmt.Println(pkg.Red(string(global.LogoContent)))
	tip()
	fmt.Println(pkg.Green("Server run at:"))
	fmt.Printf("-  Local:   %s://localhost:%d/ \r\n", "http", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/ \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Println(pkg.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/swagger/admin/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/swagger/admin/index.html \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", pkg.GetCurrentTimeStr())

	<-quit
	// Restored here, not deferred: from this point a second signal must reach
	// the default handler, so a shutdown that hangs can still be interrupted.
	disarmStopSignals()

	log.Info("Shutdown Server ... ")
	if err := shutdownServer(srv, shutdownTimeout); err != nil {
		// Not log.Fatal: that is an unconditional os.Exit(1), and Shutdown
		// reports an error exactly when connections were still in flight -
		// which is when the cleanup that follows matters most.
		log.Error("Server Shutdown: ", err)
	}
	log.Info("Server exiting")

	return nil
}

// shutdownTimeout is how long Shutdown waits for in-flight requests. It plus
// whatever cleanup follows has to stay inside the orchestrator's grace period
// - `docker stop` allows 10s by default before it sends SIGKILL.
const shutdownTimeout = 5 * time.Second

// armStopSignals registers for the stop signals and returns the channel they
// arrive on together with the function that restores the default disposition.
//
// SIGTERM is what actually arrives in production: `docker stop`, a Kubernetes
// pod deletion and `systemctl stop` all send it, and Go terminates the process
// immediately for a signal nobody listens for. Registering only os.Interrupt
// meant every graceful shutdown below the wait was dead code outside a
// terminal.
//
// Registering is separate from waiting so a caller can arm before it announces
// that it is ready: a signal that arrives between the two is delivered to the
// default handler, which for both of these means the process dies without
// running any of this.
func armStopSignals() (<-chan os.Signal, func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	return quit, func() { signal.Stop(quit) }
}

// shutdownServer stops srv, giving in-flight requests up to timeout to finish.
//
// It returns the error instead of exiting on it. A caller that exits here skips
// its own cleanup, and Shutdown fails precisely when there was something left
// to clean up after.
func shutdownServer(srv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return srv.Shutdown(ctx)
}

// runStartupHooks runs the router registries and then the before callbacks.
//
// The package-level slice runs first and in its existing order, so a fork that
// only ever appended to AppRouters sees no change at all. The core registry
// runs second, through RunAppRouters: a module can register through
// sdk.Runtime.SetAppRouters and no longer has to import this command package -
// which is a main package's plumbing - just to be routed.
//
// The loop over the core registry now lives in core, which is what brings the
// panic guard and the registration seal with it. RunBefore closes a gap rather
// than moving one: the open-source edition never executed the before callbacks
// at all, so SetBefore was accepted and silently ignored. It has to stay ahead
// of ListenAndServe, because a callback registered WithFatal exits the process
// and that must not happen to one that is already serving.
func runStartupHooks() {
	for _, f := range AppRouters {
		f()
	}
	sdk.Runtime.RunAppRouters()
	sdk.Runtime.RunBefore()
}

//var Router runtime.Router

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`go-admin `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}

func initRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		h = gin.New()
		sdk.Runtime.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		//os.Exit(-1)
	}
	if config.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	//r.Use(middleware.Metrics())
	r.Use(common.Sentinel()).
		Use(common.RequestId(pkg.TrafficKey)).
		Use(api.SetRequestLogger)

	common.InitMiddleware(r)

}
