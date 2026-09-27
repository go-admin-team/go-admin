package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	mycasbin "github.com/go-admin-team/go-admin-core/v2/casbin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
)

// Every route the code generator registers, as a request and as the policy
// that grants it.
var genRoutes = []struct{ method, path, policy string }{
	{"GET", "/api/v1/gen/preview/1", "/api/v1/gen/preview/:tableId"},
	{"GET", "/api/v1/gen/toproject/1", "/api/v1/gen/toproject/:tableId"},
	{"GET", "/api/v1/gen/apitofile/1", "/api/v1/gen/apitofile/:tableId"},
	{"GET", "/api/v1/gen/todb/1", "/api/v1/gen/todb/:tableId"},
	{"GET", "/api/v1/gen/tabletree", "/api/v1/gen/tabletree"},
	{"GET", "/api/v1/db/tables/page", "/api/v1/db/tables/page"},
	{"GET", "/api/v1/db/columns/page", "/api/v1/db/columns/page"},
	{"GET", "/api/v1/sys/tables/page", "/api/v1/sys/tables/page"},
	{"GET", "/api/v1/sys/tables/info", "/api/v1/sys/tables/info"},
	{"GET", "/api/v1/sys/tables/info/1", "/api/v1/sys/tables/info/:tableId"},
	{"POST", "/api/v1/sys/tables/info", "/api/v1/sys/tables/info"},
	{"PUT", "/api/v1/sys/tables/info", "/api/v1/sys/tables/info"},
	{"DELETE", "/api/v1/sys/tables/info/1", "/api/v1/sys/tables/info/:tableId"},
}

// genAuthzEngine serves the generator's routes behind a real JWT middleware
// and a real casbin enforcer for host, holding the policies given.
func genAuthzEngine(t *testing.T, host string, grant func(add func(role, path, method string))) (*gin.Engine, *jwt.GinJWTMiddleware) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	previousMode := config.ApplicationConfig.Mode
	config.ApplicationConfig.Mode = "dev" // the writing routes exist only here
	t.Cleanup(func() { config.ApplicationConfig.Mode = previousMode })

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	previousInterval := mycasbin.ReloadInterval
	mycasbin.ReloadInterval = 0
	t.Cleanup(func() { mycasbin.ReloadInterval = previousInterval })
	e := mycasbin.Setup(db, host)
	grant(func(role, path, method string) {
		if _, err := e.AddPolicy(role, path, method); err != nil {
			t.Fatal(err)
		}
	})
	previous := sdk.Runtime.GetCasbinByTenant(host)
	sdk.Runtime.SetCasbinByTenant(host, e)
	t.Cleanup(func() { sdk.Runtime.SetCasbinByTenant(host, previous) })

	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "test",
		Key:           []byte("test-key"),
		Timeout:       time.Hour,
		PayloadFunc:   func(data interface{}) jwt.MapClaims { return data.(jwt.MapClaims) },
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	// A request AuthCheckRole lets through reaches a handler with no logger or
	// database in its context, which panics. Recovery turns that into a 500:
	// past the check, which is all these tests ask.
	r.Use(gin.RecoveryWithWriter(io.Discard))
	v1 := r.Group("/api/v1")
	sysNoCheckRoleRouter(v1, mw)
	registerDBRouter(v1, mw)
	registerSysTableRouter(v1, mw)
	return r, mw
}

// deniedByRole reports whether AuthCheckRole refused the request. It answers
// HTTP 200 with code 403 in the body; anything else got past it, even if the
// handler behind it then failed for want of a database.
func deniedByRole(t *testing.T, r *gin.Engine, mw *jwt.GinJWTMiddleware, host, role, method, path string) bool {
	t.Helper()
	token, _, err := mw.TokenGenerator(jwt.MapClaims{"identity": 7, "rolekey": role, "roleid": 2})
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader("{}"))
	req.Host = host
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var body struct{ Code int }
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return body.Code == 403 && strings.Contains(w.Body.String(), "没有该接口访问权限")
}

func TestGeneratorRoutesRequireARoleThatHoldsThem(t *testing.T) {
	const host = "gen-authz.test"
	r, mw := genAuthzEngine(t, host, func(add func(role, path, method string)) {
		for _, route := range genRoutes {
			add("developer", route.policy, route.method)
		}
	})

	for _, route := range genRoutes {
		name := route.method + " " + route.path
		if !deniedByRole(t, r, mw, host, "clerk", route.method, route.path) {
			t.Errorf("%s: a role holding no policy got past AuthCheckRole", name)
		}
		if deniedByRole(t, r, mw, host, "developer", route.method, route.path) {
			t.Errorf("%s: a role holding the policy was refused", name)
		}
		if deniedByRole(t, r, mw, host, "admin", route.method, route.path) {
			t.Errorf("%s: admin was refused", name)
		}
	}
}

// The captcha shares a registration function with the generator's routes and
// has to stay reachable before anyone has logged in.
func TestCaptchaStaysPublic(t *testing.T) {
	r, _ := genAuthzEngine(t, "gen-authz-captcha.test", func(func(string, string, string)) {})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/captcha", nil))
	if w.Code == http.StatusUnauthorized || strings.Contains(w.Body.String(), "没有该接口访问权限") {
		t.Errorf("captcha answered %d %s", w.Code, w.Body.String())
	}
}
