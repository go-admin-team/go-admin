package middleware

import (
	"time"

	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"go-admin/common/middleware/handler"
)

// authMiddleware is the single JWT middleware instance the whole process
// shares. InitMiddleware builds it once, before any module registers its
// routes; GetAuthMiddleware is how a module gets it back instead of calling
// AuthInit itself and building another, functionally-equivalent-but-distinct
// instance.
var authMiddleware *jwt.GinJWTMiddleware

// AuthInit jwt验证new
func AuthInit() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		timeout = time.Duration(876010) * time.Hour
	} else {
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.JwtConfig.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     handler.PayloadFunc,
		IdentityHandler: handler.IdentityHandler,
		Authenticator:   handler.Authenticator,
		Authorizator:    handler.Authorizator,
		Unauthorized:    handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}

// GetAuthMiddleware returns the shared JWT middleware instance InitMiddleware
// built at startup. Application modules (app/admin, app/jobs, app/other,
// app/demo) call this instead of AuthInit so their router chains - which
// still need the instance itself for authMiddleware.MiddlewareFunc() and
// authMiddleware.LoginHandler, not just the bound closure registered under
// sdk.Runtime's JwtTokenCheck key - end up using the same instance the host
// registered, rather than one each.
//
// It fails loudly instead of returning nil: an InitRouter that runs before
// InitMiddleware has a real startup-ordering bug, not a case to paper over
// with a nil *jwt.GinJWTMiddleware that would panic much further down the
// call chain with a far less useful stack trace.
func GetAuthMiddleware() *jwt.GinJWTMiddleware {
	if authMiddleware == nil {
		log.Fatal("JWT middleware not initialized; InitMiddleware must run before any module's InitRouter")
	}
	return authMiddleware
}
