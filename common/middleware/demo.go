package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// defaultDemoMsg is what a refused request is told when nothing is configured.
//
// It is the string this middleware used to carry hard-coded, kept verbatim so
// that a deployment which never set application.demomsg is answered exactly as
// it was before.
const defaultDemoMsg = "谢谢您的参与，但为了大家更好的体验，所以本次提交就算了吧！\U0001F600\U0001F600\U0001F600"

// demoWriteRoutes are routes that change something despite being registered as
// GET, so the method alone does not say whether they are safe to serve.
//
// All three belong to the code generator: two write Go source files onto the
// server's filesystem and the third inserts menus, APIs and casbin rules into
// the database. They are registered under a group whose own name says it does
// no role check, and a demo deployment lets anybody log in - so on a demo host
// they were reachable by any visitor, and the menus one had in fact been used.
//
// Spelled as gin route patterns, which is what Context.FullPath returns, so a
// path parameter matches whatever value it is given.
//
// This list cannot be checked from here: common/ may not import app/, so this
// package cannot see which routes exist. What keeps it honest is a test beside
// the routes themselves - see app/other/router - which registers them and
// fails if any entry here has stopped being a real route.
//
// It also does not close the general hole. Nothing stops the next GET handler
// that writes something from being added without an entry here, and no static
// check can tell a handler that writes from one that reads. Demo mode refuses
// the routes it has been told about; that is the whole of the guarantee.
var demoWriteRoutes = map[string]bool{
	"/api/v1/gen/toproject/:tableId": true,
	"/api/v1/gen/apitofile/:tableId": true,
	"/api/v1/gen/todb/:tableId":      true,
}

// DemoWriteRoutes returns the routes demo mode refuses despite their method.
//
// Exported only so the test that lives beside the route registrations can
// check every one of them still exists; nothing else should need it.
func DemoWriteRoutes() []string {
	out := make([]string, 0, len(demoWriteRoutes))
	for route := range demoWriteRoutes {
		out = append(out, route)
	}
	return out
}

// demoAllows reports whether demo mode lets a request through.
//
// route is the matched gin route pattern and uri the raw request target; the
// two are different things and both are needed. The route is what identifies a
// handler regardless of the values in its path parameters, and it is empty for
// a request that matched nothing - which is why the login and logout checks
// still read the raw target, as they always did.
func demoAllows(method, route, uri string) bool {
	if demoWriteRoutes[route] {
		return false
	}
	return method == http.MethodGet ||
		method == http.MethodOptions ||
		uri == "/api/v1/login" ||
		uri == "/api/v1/logout"
}

// demoMessage is the answer a refused request gets.
//
// application.demomsg has been in the configuration all along and nothing read
// it: the message was hard-coded here, and the demo host's configured string
// happened to be identical, so the setting looked like it worked. An empty
// value falls back rather than answering with nothing.
func demoMessage() string {
	if msg := config.ApplicationConfig.DemoMsg; msg != "" {
		return msg
	}
	return defaultDemoMsg
}

// DemoEvn refuses anything that would change state while mode is demo.
func DemoEvn() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ApplicationConfig.Mode != "demo" {
			c.Next()
			return
		}
		if demoAllows(c.Request.Method, c.FullPath(), c.Request.RequestURI) {
			c.Next()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  demoMessage(),
		})
		c.Abort()
	}
}
