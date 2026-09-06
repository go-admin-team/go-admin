package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// demoMode puts the process in demo mode for one test and puts it back.
func demoMode(t *testing.T, mode, msg string) {
	t.Helper()
	previousMode, previousMsg := config.ApplicationConfig.Mode, config.ApplicationConfig.DemoMsg
	t.Cleanup(func() {
		config.ApplicationConfig.Mode = previousMode
		config.ApplicationConfig.DemoMsg = previousMsg
	})
	config.ApplicationConfig.Mode, config.ApplicationConfig.DemoMsg = mode, msg
}

// The method is not enough on its own. Three of the generator's routes are
// registered as GET and write anyway - two of them onto the server's
// filesystem, one into the database - so a guard that reads only the method
// serves them to anybody who can log in, which on a demo host is everybody.
func TestDemoRefusesTheWritesThatAreServedOverGET(t *testing.T) {
	const login = "/api/v1/login"
	for _, tc := range []struct {
		name        string
		method      string
		route, uri  string
		wantThrough bool
	}{
		{"a plain read", http.MethodGet, "/api/v1/dept", "/api/v1/dept", true},
		{"a write, by method", http.MethodPost, "/api/v1/dept", "/api/v1/dept", false},
		{"login is how a visitor gets in", http.MethodPost, login, login, true},
		{"logout", http.MethodPost, "/api/v1/logout", "/api/v1/logout", true},
		{"preflight", http.MethodOptions, "/api/v1/dept", "/api/v1/dept", true},
		// A request that matched no route has an empty pattern, and the guard
		// still has to refuse it by method - this is what a POST to a path
		// that does not exist looks like from in here.
		{"a write to nothing at all", http.MethodPost, "", "/api/v1/__probe__", false},

		// The three this change is about.
		{"generator writes the database", http.MethodGet,
			"/api/v1/gen/todb/:tableId", "/api/v1/gen/todb/3", false},
		{"generator writes source files", http.MethodGet,
			"/api/v1/gen/toproject/:tableId", "/api/v1/gen/toproject/3", false},
		{"generator writes an api file", http.MethodGet,
			"/api/v1/gen/apitofile/:tableId", "/api/v1/gen/apitofile/3", false},

		// The read-only half of the generator has to keep working, or the demo
		// host cannot demonstrate the feature at all. Refusing too much is as
		// much of a defect as refusing too little.
		{"generator preview stays available", http.MethodGet,
			"/api/v1/gen/preview/:tableId", "/api/v1/gen/preview/3", true},
		{"generator table tree stays available", http.MethodGet,
			"/api/v1/gen/tabletree", "/api/v1/gen/tabletree", true},
		{"table list stays available", http.MethodGet,
			"/api/v1/db/tables/page", "/api/v1/db/tables/page", true},
		{"column list stays available", http.MethodGet,
			"/api/v1/db/columns/page", "/api/v1/db/columns/page", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := demoAllows(tc.method, tc.route, tc.uri); got != tc.wantThrough {
				t.Errorf("demoAllows(%s %s) = %v, want %v", tc.method, tc.route, got, tc.wantThrough)
			}
		})
	}
}

// Everything above is about demo mode only. A deployment that is not a demo
// runs the generator for real, and a guard that reached it there would have
// taken the feature away from every production install.
func TestOutsideDemoModeNothingIsRefused(t *testing.T) {
	demoMode(t, "prod", "")
	gin.SetMode(gin.TestMode)

	for _, route := range append(DemoWriteRoutes(), "/api/v1/dept") {
		t.Run(route, func(t *testing.T) {
			r := gin.New()
			r.Use(DemoEvn())
			r.GET(route, func(c *gin.Context) { c.String(http.StatusOK, "served") })

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, requestFor(route), nil))
			if w.Body.String() != "served" {
				t.Errorf("answered %q; outside demo mode the handler must run", w.Body.String())
			}
		})
	}
}

// The refusal has to come back as the demo message rather than a 403 or a 404:
// the front end shows it to the visitor, and the point of a demo host is that
// being turned away is explained.
func TestDemoRefusalCarriesTheConfiguredMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const route = "/api/v1/gen/todb/:tableId"

	for _, tc := range []struct {
		name, configured, want string
	}{
		{"configured", "come back tomorrow", "come back tomorrow"},
		// A deployment that never set application.demomsg keeps the answer it
		// already had; an empty setting must not become an empty message.
		{"not configured", "", defaultDemoMsg},
	} {
		t.Run(tc.name, func(t *testing.T) {
			demoMode(t, "demo", tc.configured)

			r := gin.New()
			r.Use(DemoEvn())
			r.GET(route, func(c *gin.Context) { c.String(http.StatusOK, "served") })

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, requestFor(route), nil))

			if w.Code != http.StatusOK {
				t.Errorf("answered %d, want 200 so the front end reads the body", w.Code)
			}
			if body := w.Body.String(); !strings.Contains(body, tc.want) {
				t.Errorf("body %q does not carry %q", body, tc.want)
			}
			if strings.Contains(w.Body.String(), "served") {
				t.Error("the handler ran; the request was supposed to be refused")
			}
		})
	}
}

// requestFor turns a route pattern into a request target by giving every path
// parameter a value.
func requestFor(route string) string {
	segments := strings.Split(route, "/")
	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			segments[i] = "1"
		}
	}
	return strings.Join(segments, "/")
}
