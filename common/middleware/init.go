package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	"go-admin/common/actions"
)

// These alias core's own constants (see sdk/runtime.GetHandlerFunc's contract
// doc, section 9) rather than redeclaring the same three strings, so a typo
// here can no longer split registration and lookup into two different keys
// that both happen to compile.
const (
	JwtTokenCheck   = runtime.JwtTokenCheck
	RoleCheck       = runtime.RoleCheck
	PermissionCheck = runtime.PermissionCheck
)

func InitMiddleware(r *gin.Engine) {
	r.Use(DemoEvn())
	// 数据库链接
	r.Use(WithContextDb)
	// 日志处理
	r.Use(LoggerToFile())
	// 自定义错误处理
	r.Use(CustomError)
	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// 跨域处理
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	// 链路追踪
	//r.Use(middleware.Trace())

	// Build the shared JWT middleware instance here, before any module
	// registers routes (initRouter runs ahead of runStartupHooks, which is
	// what invokes each module's InitRouter - see cmd/api/server.go). Doing
	// it once here, instead of once per module via AuthInit, is what makes
	// GetAuthMiddleware and sdk.Runtime.GetHandlerFunc(JwtTokenCheck) both
	// resolve to a single, meaningful instance instead of "whichever module
	// happened to initialize last".
	//
	// SetMiddleware must be given a bound closure (authMiddleware.MiddlewareFunc()),
	// not the unbound method expression (*jwt.GinJWTMiddleware).MiddlewareFunc:
	// the latter has no receiver bound to it, so GetHandlerFunc's type
	// assertion to gin.HandlerFunc always fails for it.
	var err error
	authMiddleware, err = AuthInit()
	if err != nil {
		// A process with no JWT middleware must not start serving requests.
		log.Fatalf("JWT Init Error, %s", err.Error())
	}
	sdk.Runtime.SetMiddleware(JwtTokenCheck, authMiddleware.MiddlewareFunc())
	sdk.Runtime.SetMiddleware(RoleCheck, AuthCheckRole())
	sdk.Runtime.SetMiddleware(PermissionCheck, actions.PermissionAction())
}
