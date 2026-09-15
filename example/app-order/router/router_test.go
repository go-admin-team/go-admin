package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	coreruntime "github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
)

func testAuthMiddleware(t *testing.T) *jwt.GinJWTMiddleware {
	t.Helper()
	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test",
		Key:              []byte("test-signing-key"),
		SigningAlgorithm: "HS256",
		Timeout:          0,
		TokenLookup:      "header: Authorization",
		TokenHeadName:    "Bearer",
	})
	if err != nil {
		t.Fatalf("building a test JWT middleware: %v", err)
	}
	return mw
}

// sdk.Runtime is a single process-wide instance (see its doc comment) with no
// way to unregister a middleware key, so this test needs RoleCheck to be
// unset - which makes it look order-dependent. It is not: the test that does
// register RoleCheck puts it back in a t.Cleanup, and the guard below turns a
// wrong order into a loud failure rather than a silent pass. Verified with
// `go test -shuffle=<seed>` on seeds that run the two in either order.
func TestRegisterRouterPanicsWithoutHostRoleCheck(t *testing.T) {
	if _, ok := sdk.Runtime.GetHandlerFunc(coreruntime.RoleCheck); ok {
		t.Fatal("RoleCheck is already registered; this test must run before any test that registers it")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("RegisterRouter did not panic with no host RoleCheck middleware registered")
		}
	}()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	RegisterRouter(v1, testAuthMiddleware(t))
}

func TestRegisterRouterMountsRoutesOnceRoleCheckIsRegistered(t *testing.T) {
	sdk.Runtime.SetMiddleware(coreruntime.RoleCheck, gin.HandlerFunc(func(c *gin.Context) { c.Next() }))
	// sdk.Runtime has no way to unregister a middleware key (SetMiddleware
	// only ever adds or overwrites - see its doc comment), so restore the
	// "as far as GetHandlerFunc is concerned, unregistered" state other
	// tests in this package depend on: a nil interface{} fails
	// GetHandlerFunc's gin.HandlerFunc type assertion the same way a never-
	// set key does. Needed for `go test -count=2` and similar re-runs
	// within one process, not for a single run.
	t.Cleanup(func() { sdk.Runtime.SetMiddleware(coreruntime.RoleCheck, nil) })

	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	RegisterRouter(v1, testAuthMiddleware(t))

	// A route that exists returns something other than 404, even if the
	// JWT/Casbin/PermissionAction chain in front of it then rejects the
	// unauthenticated test request - proving RegisterRouter actually wired
	// the route up is the point, not exercising the auth chain itself.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/order", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Errorf("GET /api/v1/order was not registered (404)")
	}
}
