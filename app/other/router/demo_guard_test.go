package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"

	"go-admin/common/middleware"
)

// registeredRoutes builds the generator's routes on an engine of its own and
// reports the patterns they were registered under.
//
// The JWT middleware is a zero value. MiddlewareFunc only closes over the
// receiver and is never called here - no request is served, the engine is
// asked what it has - so nothing dereferences it.
func registeredRoutes(t *testing.T) map[string]bool {
	t.Helper()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	v1 := r.Group("/api/v1")
	sysNoCheckRoleRouter(v1, &jwt.GinJWTMiddleware{})
	registerDBRouter(v1, &jwt.GinJWTMiddleware{})

	out := map[string]bool{}
	for _, route := range r.Routes() {
		out[route.Path] = true
	}
	return out
}

// The list of routes demo mode refuses lives in common/middleware, which may
// not import app/ and therefore cannot see whether any of them is still a
// route. This is the half that can be checked, and it is checked here because
// this is where the routes are declared: rename one, and the entry over there
// stops matching anything, demo mode silently starts serving it again, and
// nothing else would say so.
func TestEveryRouteDemoModeRefusesStillExists(t *testing.T) {
	routes := registeredRoutes(t)
	for _, guarded := range middleware.DemoWriteRoutes() {
		if !routes[guarded] {
			t.Errorf("demo mode refuses %q, but no route is registered under that pattern - "+
				"either it was renamed, or it moved to another file; the guard now matches nothing",
				guarded)
		}
	}
}

// The other direction, and the one the demo host cares about: the generator's
// read-only routes have to stay reachable, or a demo deployment cannot show
// the feature at all. Refusing too much is as much of a defect as refusing too
// little.
func TestTheGeneratorsReadOnlyRoutesAreNotRefused(t *testing.T) {
	refused := map[string]bool{}
	for _, guarded := range middleware.DemoWriteRoutes() {
		refused[guarded] = true
	}

	for _, readOnly := range []string{
		"/api/v1/gen/preview/:tableId",
		"/api/v1/gen/tabletree",
		"/api/v1/db/tables/page",
		"/api/v1/db/columns/page",
	} {
		if !registeredRoutes(t)[readOnly] {
			t.Fatalf("%s is not registered, so this test is asserting against nothing", readOnly)
		}
		if refused[readOnly] {
			t.Errorf("demo mode refuses %s, which only reads - the demo host needs it to "+
				"demonstrate the generator", readOnly)
		}
	}
}
