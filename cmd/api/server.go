package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/config/source/file"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/bootstrap"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
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

	// Registered before the configuration is read. SetupConfig announces
	// AfterResource as soon as the callbacks that build the resources have
	// run, so a hook added after that call would miss the first round and the
	// queue would have no consumers until somebody edited the config file.
	sdk.Runtime.SetPhase(runtime.AfterResource, attachQueueConsumers)

	// On AfterListen rather than on a bare goroutine from run(). Two reasons:
	// the phase runs behind core's panic guard, which does not reach across a
	// goroutine boundary - a panic while loading jobs used to take the whole
	// process down with a stack that named this file - and the jobs it starts
	// can call the API, which is only true once the socket is accepting.
	sdk.Runtime.SetPhase(runtime.AfterListen, startCronJobs)

	//1. 读取配置
	bootstrap.SetupConfig(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		storage.Setup,
	)

	usageStr := `starting api server...`
	log.Info(usageStr)
}

// startCronJobs registers the job implementations and starts a scheduler for
// every tenant database.
//
// It is synchronous, like the phase that runs it. jobs.Setup returns now that
// the `select {}` at the end of its per-tenant setup is gone, which is what
// makes that possible; while it was there this could only be a goroutine, and
// a goroutine is outside the panic guard.
func startCronJobs() {
	jobs.InitJob()
	jobs.Setup(sdk.Runtime.GetAllDb())
}

// attachedQueue is the queue generation the consumers are attached to, plus
// one, so that the zero value means "attached to nothing yet". Written from
// the goroutine running the phase, read from the next one - rounds never
// overlap, but they are not the same goroutine.
var attachedQueue atomic.Uint64

// attachQueueConsumers registers the log consumers against the queue that is
// current, and starts it.
//
// It runs on AfterResource, so it runs again after every configuration reload
// - and it has to. A reload rebuilds the queue adapter, and consumers
// registered against the one that existed at start-up are attached to an
// adapter nobody publishes to any more, so the login and operation logs stop
// being written with nothing said about it.
//
// It is therefore idempotent with respect to a given queue rather than "does
// nothing the second time": a new adapter gets a fresh set of consumers, the
// same one gets none. Registering twice on the same queue would give every
// message two consumers and write every log row twice.
//
// Generation 0 means the configuration has no queue section, so nothing was
// installed and GetQueuePrefix hands back the runtime's own memory queue.
// That case still gets consumers - it is what the previous unconditional
// registration did, and dropping it would silently stop logging for anyone who
// commented the section out - it just never gets them twice.
func attachQueueConsumers() {
	attachConsumersOnce(storage.QueueGeneration(), sdk.Runtime.GetQueuePrefix(""))
}

// attachConsumersOnce puts the log consumers on q and starts it, unless gen
// says this queue already has them.
//
// Split out from attachQueueConsumers so that the order and the once-ness can
// be checked against a queue the test controls: the sequence that matters here
// cannot be read back out of a real adapter.
func attachConsumersOnce(gen uint64, q corestorage.AdapterQueue) {
	if attachedQueue.Load() == gen+1 {
		return
	}
	attachedQueue.Store(gen + 1)

	//注册监听函数
	q.Register(global.LoginLog, models.SaveLoginLog)
	q.Register(global.OperateLog, models.SaveOperaLog)
	q.Register(global.ApiCheck, models.SaveSysApi)

	// Started only now, and by whoever registered. setupQueue deliberately
	// leaves it stopped: a queue that is already running refuses further
	// registration, and the adapter in this path drops that error on the
	// floor, so starting first loses consumers without a word.
	go q.Run()
}

func run() error {
	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	buildRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler:      sdk.Runtime.GetEngine(),
		ReadTimeout:  time.Duration(config.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ApplicationConfig.WriterTimeout) * time.Second,
	}

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

	if err := startServing(srv, config.SslConfig.Enable, config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil {
		return err
	}

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

	// Said before anything is taken apart. A configuration reload arriving in
	// this window would otherwise re-run AfterResource - rebuilding the pool
	// and the queue adapter, and re-registering consumers - on top of cleanup
	// that has already run.
	sdk.Runtime.BeginShutdown()

	log.Info("Shutdown Server ... ")
	if err := shutdownServer(srv, shutdownTimeout); err != nil {
		// Not log.Fatal: that is an unconditional os.Exit(1), and Shutdown
		// reports an error exactly when connections were still in flight -
		// which is when the cleanup that follows matters most.
		log.Error("Server Shutdown: ", err)
	}

	// Runs whether or not the line above reported an error, for that reason.
	if err := runShutdownHooks(cleanupTimeout); err != nil {
		log.Error("Cleanup: ", err)
	}
	log.Info("Server exiting")

	return nil
}

// shutdownTimeout is how long Shutdown waits for in-flight requests, and
// cleanupTimeout how long the BeforeExit callbacks get after it.
//
// They are consumed one after the other, so the two together are what has to
// stay inside the orchestrator's grace period: `docker stop` allows 10s by
// default before it sends SIGKILL, and 5+3 leaves room for the process to
// finish returning. Raising either without lowering the other buys nothing -
// the budget that runs out is the orchestrator's.
const (
	shutdownTimeout = 5 * time.Second
	cleanupTimeout  = 3 * time.Second
)

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

// startServing binds srv.Addr, hands the listener to srv on its own goroutine,
// and announces AfterListen.
//
// The bind is done here rather than left to ListenAndServe, which binds on the
// goroutine that serves. That put the failure every deployment actually hits -
// "address already in use" - on a goroutine nobody was reading, so the banner
// went on to claim the server was up, and there would be no way to keep
// AfterListen from announcing a socket that does not exist. A hook there is
// promised a reachable port; the only way to keep that promise is for the bind
// to have already happened on this goroutine.
//
// AfterListen is announced synchronously. Running it in a goroutine to save the
// few milliseconds would let it overlap the shutdown: on a fast SIGTERM the
// cleanup callbacks could finish before the startup ones had.
//
// Both ways of failing to start are therefore checked before the announcement:
// the bind, and - with ssl enabled - the certificate.
func startServing(srv *http.Server, useTLS bool, pem, key string) error {
	if useTLS {
		// Read the certificate before anything is announced. ServeTLS reads
		// these files itself, but on the serving goroutine - so a bad
		// certificate used to surface after AfterListen had already promised a
		// reachable port. Loading it here costs one extra read and moves the
		// failure onto this goroutine, where run() can return it.
		//
		// ServeTLS still does the real work below rather than this handing it a
		// tls.Listener: that is what sets up HTTP/2 negotiation, and taking it
		// over here would quietly drop h2 for every TLS deployment.
		if _, err := tls.LoadX509KeyPair(pem, key); err != nil {
			return errors.Wrap(err, "tls certificate")
		}
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	go func() {
		// 服务连接
		var err error
		if useTLS {
			err = srv.ServeTLS(ln, pem, key)
		} else {
			err = srv.Serve(ln)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Still fatal, as it was. Neither the bind nor the certificate is
			// among the errors that reach here any more - both are checked
			// above, on the caller's goroutine. What is left is a serve that
			// failed after the port was taken, and carrying on would park the
			// process on <-quit with nothing serving.
			log.Fatal("serve: ", err)
		}
	}()

	sdk.Runtime.RunPhase(runtime.AfterListen)
	return nil
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

// runShutdownHooks runs the BeforeExit callbacks with timeout to share.
//
// What the budget bounds is the wait, not the work. When it is gone RunShutdown
// stops waiting and returns; a callback that never looks at its context carries
// on until the process exits, and may leave a partial write behind. Go cannot
// cancel a function that does not check for cancellation, which is why the
// callbacks are handed a context at all.
func runShutdownHooks(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return sdk.Runtime.RunShutdown(ctx)
}

// buildRouter announces BeforeRouter, builds the engine, and then drains the
// startup registries.
//
// The order is the contract. BeforeRouter is the last point at which a module
// can still affect how routes are built, so it has to run while there is no
// engine yet. The before registry runStartupHooks drains is a different moment
// despite the name: those callbacks run after initRouter has built the engine.
// Two lines apart, and describing them as equivalent is a mistake this
// repository has already made once in writing.
func buildRouter() {
	sdk.Runtime.RunPhase(runtime.BeforeRouter)
	initRouter()
	runStartupHooks()
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
