package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"

	"go-admin/common/middleware"
)

// registeredRoutes builds the generator's routes on an engine of its own and
// reports the patterns they were registered under.
//
// The mode has to be given rather than inherited, because it now decides what
// gets registered: a test that leaves it at the zero value would be asking
// about a mode no deployment runs in, and would pass whether or not the gate
// works.
//
// It is put back before this returns, not at the end of the test. t.Cleanup
// would leave the mode set for everything the caller does afterwards, so a
// caller that went on to assert something mode-dependent would be reading a
// value this helper left behind rather than one it chose. A caller that does
// want the mode set has to set it, which is visible where it happens.
//
// The JWT middleware is a zero value. MiddlewareFunc only closes over the
// receiver and is never called here - no request is served, the engine is
// asked what it has - so nothing dereferences it.
func registeredRoutes(t *testing.T, mode string) map[string]bool {
	t.Helper()
	gin.SetMode(gin.TestMode)

	previous := config.ApplicationConfig.Mode
	defer func() { config.ApplicationConfig.Mode = previous }()
	config.ApplicationConfig.Mode = mode

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
	routes := registeredRoutes(t, "demo")
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
		if !registeredRoutes(t, "demo")[readOnly] {
			t.Fatalf("%s is not registered, so this test is asserting against nothing", readOnly)
		}
		if refused[readOnly] {
			t.Errorf("demo mode refuses %s, which only reads - the demo host needs it to "+
				"demonstrate the generator", readOnly)
		}
	}
}
