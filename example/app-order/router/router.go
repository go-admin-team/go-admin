// Package router wires app-order's four routes onto a host's gin engine.
package router

import (
	"github.com/gin-gonic/gin"

	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"
	coreruntime "github.com/go-admin-team/go-admin-core/v2/sdk/runtime"

	"github.com/go-admin-team/example-app-order/apis"
)

// RegisterRouter mounts app-order's routes under v1.
//
// Its signature - (v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware)
// - is not app-order's own invention: it is the exact shape every in-tree
// go-admin app router package already registers into its own routerCheckRole
// slice (see app/demo/router/demo_product.go), so a host installs this
// exactly where it installs its own app/*/router packages: one file under
// cmd/api/ that imports this package and appends RegisterRouter (adjusted to
// the host's own registration slice's calling convention) - see
// cmd/api/demo.go for the pattern.
//
// authMiddleware is taken as an explicit parameter rather than fetched
// through sdk.Runtime.GetHandlerFunc(coreruntime.JwtTokenCheck). As of this
// writing the reference host (go-admin's common/middleware/init.go) registers
// that key with an unbound method expression -
// sdk.Runtime.SetMiddleware(JwtTokenCheck, (*jwt.GinJWTMiddleware).MiddlewareFunc)
// - which is exactly the shape GetHandlerFunc's own doc comment warns
// against: the stored value's type is func(*jwt.GinJWTMiddleware)
// gin.HandlerFunc, not gin.HandlerFunc, so GetHandlerFunc's type assertion
// fails and it reports ok=false every time, for every caller, not just this
// one. Taking authMiddleware directly sidesteps that live bug and matches
// what every in-tree app already does.
func RegisterRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	roleCheck, ok := sdk.Runtime.GetHandlerFunc(coreruntime.RoleCheck)
	if !ok {
		// A host that has not wired up RoleCheck has not wired up Casbin
		// authorization at all. Registering these routes without it would
		// silently serve every order to every authenticated caller
		// regardless of role - fail loud at startup instead, the same way
		// PermissionAction fails loud (Abort, not c.Next) when its own
		// database lookup errors. See contract/actions.PermissionAction's
		// doc comment for the same reasoning applied to data-scope instead
		// of role.
		panic("app-order: host has not registered core's " + coreruntime.RoleCheck +
			" middleware (sdk.Runtime.SetMiddleware); refusing to mount unauthorized order routes")
	}

	e := apis.Order{}
	r := v1.Group("/order").
		Use(authMiddleware.MiddlewareFunc()).
		Use(roleCheck)
	{
		// actions.PermissionAction is imported directly from core - a plain
		// function, not something fetched through sdk.Runtime - because
		// unlike RoleCheck's Casbin policy tables (host-owned; see
		// contract/actions's package doc), the data-scope machinery it
		// installs has no host-specific state at all.
		r.GET("", actions.PermissionAction(), e.GetPage)
		r.GET("/:id", actions.PermissionAction(), e.Get)
		r.POST("", e.Create)
		r.PUT("/:id/pay", actions.PermissionAction(), e.Pay)
	}
}
