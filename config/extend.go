package config

import (
	"fmt"
	"strings"
)

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap      AMap // 这里配置对应配置文件的结构即可
	FileStore FileStore
	RateLimit RateLimit
	Shutdown  Shutdown
}

// DefaultInboundQPS is the limit applied when nothing is configured. It is the
// value that used to be hard-coded in the middleware, so an existing deployment
// that adds nothing to settings.yml keeps the behaviour it already had.
const DefaultInboundQPS = 200

// RateLimit 全局入站限流。
//
//	extend:
//	  ratelimit:
//	    inboundqps: 200   # 每秒入站请求上限；填 0 关闭限流
//
// The threshold used to live in common/middleware/sentinel.go as a constant,
// which made 200 QPS the ceiling of every deployment with nothing in the
// configuration to reveal it.
type RateLimit struct {
	// InboundQPS caps inbound requests per second across the process.
	//
	// Absent means DefaultInboundQPS, zero disables the limiter, and a positive
	// value is the threshold. The pointer is what separates "not configured"
	// from "configured to zero" - the two need different answers and a plain
	// float64 cannot tell them apart.
	InboundQPS *float64
}

// Threshold reports the limit to apply. Zero means no limiting.
func (r RateLimit) Threshold() float64 {
	if r.InboundQPS == nil {
		return DefaultInboundQPS
	}
	return *r.InboundQPS
}

type AMap struct {
	Key string
}

// FileStore 对象存储。上传接口的 source 参数决定走哪一家：2 是阿里云，3 是七牛。
// 没有填的那一家在被请求时返回明确错误，而不是上传到别处或者崩溃。
//
// common/file_store 里还实现了华为云 OBS，但上传接口没有对应的 source 取值，
// 所以这里也不为它提供配置。
type FileStore struct {
	AliYun ObjectStore
	QiNiu  ObjectStore
}

type ObjectStore struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// Configured reports whether enough was filled in to attempt a connection.
func (o ObjectStore) Configured() bool {
	return o.Endpoint != "" && o.AccessKeyID != "" && o.AccessKeySecret != "" && o.BucketName != ""
}

// Default budgets for a graceful shutdown, in seconds. Each applies to the
// matching field of extend.shutdown when that field is absent, and together
// they are what the process spent before the section existed - so a deployment
// that configures nothing keeps the shutdown it already had.
//
// The drain default is zero deliberately. The three budgets are spent one
// after the other, and once their sum reaches the orchestrator's stop grace
// period the process is killed part-way through its cleanup callbacks, which
// is worse than not draining at all. `docker stop` allows ten seconds by
// default and 5+3 already leaves little room, so a non-zero default here would
// slow down every existing shutdown to buy something only a load balancer that
// polls /ready can collect.
const (
	DefaultDrainSeconds   = 0
	DefaultServerSeconds  = 5
	DefaultCleanupSeconds = 3
)

// Shutdown is how long a graceful shutdown may spend, stage by stage.
//
//	extend:
//	  shutdown:
//	    drain: 0
//	    server: 5
//	    cleanup: 3
//	    grace: 30
//
// Every field is a pointer for the reason RateLimit.InboundQPS is: nil means
// "not configured" and takes the default, while a value that was written down
// is taken literally, zero included. Without that separation `server: 0` - do
// not wait for in-flight requests, which is a reasonable thing to ask under a
// very short grace period - could not be said at all, and `drain: 0` would
// have to mean something different from `server: 0` in the same section.
type Shutdown struct {
	// Drain is how long to keep serving normally after a stop signal arrives.
	// Throughout it /ready answers 503 and keep-alive is switched off, which
	// is what gives whatever routes traffic here time to stop routing it
	// before the listener closes. Zero is no window: the readiness flip and
	// the listener closing are then microseconds apart and nothing observes
	// the first.
	//
	// What the window is worth depends on who does the removing and on what
	// basis; the package comment in common/health has the two cases, and they
	// do not want the same value.
	Drain *int
	// Server is how long the server waits for in-flight requests once the
	// listener is closed.
	Server *int
	// Cleanup is how long the BeforeExit callbacks get after that.
	Cleanup *int
	// Grace is the stop grace period the orchestrator gives this process -
	// `docker stop --timeout`, or terminationGracePeriodSeconds. Nothing reads
	// it during a shutdown; it exists so start-up can say whether the budget
	// fits inside it. Absent means no comparison is made, because the
	// reference values differ threefold between runtimes and a fixed threshold
	// would warn about configurations that are correct.
	Grace *int
}

// ShutdownBudget is what a shutdown will actually spend, in seconds, after the
// fallbacks have been applied.
type ShutdownBudget struct {
	Drain   int
	Server  int
	Cleanup int
	// Grace is zero when extend.shutdown.grace was not configured.
	Grace int
}

// Budget resolves the configured section into the values that will be spent.
//
// A negative is refused rather than corrected. A wait cannot be negative, so
// there is no reading of one to honour, and quietly turning it into zero would
// be the failure this whole section exists to remove: written down, accepted,
// and not what happens. It is returned as an error rather than reported here
// so that the rule can be checked without ending the process.
func (s Shutdown) Budget() (ShutdownBudget, error) {
	var negative []string
	for _, f := range []struct {
		name  string
		value *int
	}{
		{"drain", s.Drain},
		{"server", s.Server},
		{"cleanup", s.Cleanup},
		{"grace", s.Grace},
	} {
		if f.value != nil && *f.value < 0 {
			negative = append(negative, fmt.Sprintf("%s: %d", f.name, *f.value))
		}
	}
	if len(negative) > 0 {
		return ShutdownBudget{}, fmt.Errorf(
			"extend.shutdown was given a negative number of seconds (%s); "+
				"a wait cannot be negative, and 0 is how to say \"do not wait\"",
			strings.Join(negative, ", "))
	}

	return ShutdownBudget{
		Drain:   budgetSeconds(s.Drain, DefaultDrainSeconds),
		Server:  budgetSeconds(s.Server, DefaultServerSeconds),
		Cleanup: budgetSeconds(s.Cleanup, DefaultCleanupSeconds),
		Grace:   budgetSeconds(s.Grace, 0),
	}, nil
}

func budgetSeconds(configured *int, fallback int) int {
	if configured != nil {
		return *configured
	}
	return fallback
}

// Total is the whole of the shutdown, since the three stages run one after the
// other.
func (b ShutdownBudget) Total() int { return b.Drain + b.Server + b.Cleanup }

// Overrun reports how many seconds have to be found for the budget to fit
// inside the configured grace period. It is zero when no grace period was
// configured and when the budget already fits.
//
// Fitting means strictly less: the grace period is when SIGKILL is sent, so a
// budget that ends exactly then leaves the last callback no time to return.
func (b ShutdownBudget) Overrun() int {
	if b.Grace <= 0 || b.Total() < b.Grace {
		return 0
	}
	return b.Total() - b.Grace + 1
}
