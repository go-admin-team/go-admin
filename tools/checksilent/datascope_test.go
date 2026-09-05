package main

import (
	"strings"
	"testing"
)

// apisFile is a handler package with two methods: one that reads the caller's
// data permission and one that does not.
const apisFile = `package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/actions"
)

type SysUser struct{}

func (e SysUser) Scoped(c *gin.Context) {
	p := actions.GetPermissionFromContext(c)
	_ = p
}

func (e SysUser) Unscoped(c *gin.Context) {}
`

func routerFile(uses string) string {
	return `package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
)

var _ = actions.PermissionAction

func register(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("/sys-user")` + uses + `
	{
		r.GET("/:id", api.Scoped)
	}
}
`
}

// The mistake itself: a handler that reads the permission, on a group that
// never installs the middleware which puts one there.
func TestDataScopeRouteWithoutTheMiddlewareIsReported(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/apis/sys_user.go":   apisFile,
		"app/admin/router/sys_user.go": routerFile(`.Use(gin.Logger())`),
	})
	f := requireOne(t, check(t, root, options{}), checkDataScopeRoute)
	for _, want := range []string{`GET "/sys-user/:id"`, "SysUser.Scoped", "PermissionAction"} {
		if !strings.Contains(f.Message, want) {
			t.Errorf("message does not mention %q:\n%s", want, f.Message)
		}
	}
}

// The middleware installed in the chain is the fix, and must silence it.
func TestDataScopeRouteWithTheMiddlewareIsQuiet(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/apis/sys_user.go":   apisFile,
		"app/admin/router/sys_user.go": routerFile(`.Use(gin.Logger()).Use(actions.PermissionAction())`),
	})
	if got := only(t, check(t, root, options{}), checkDataScopeRoute); len(got) != 0 {
		t.Errorf("reported %d findings for a guarded group:\n%v", len(got), got)
	}
}

// The other fix - the handler stops reading the permission - must silence it
// too. Reporting a route whose handler needs no scope would push people to
// install middleware they do not want, which is how /getinfo would have been
// "fixed" into rejecting every user who did not create their own account.
func TestARouteWhoseHandlerReadsNoPermissionIsQuiet(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/apis/sys_user.go": apisFile,
		"app/admin/router/sys_user.go": `package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
)

func register(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("")
	{
		r.GET("/getinfo", api.Unscoped)
	}
}
`,
	})
	if got := only(t, check(t, root, options{}), checkDataScopeRoute); len(got) != 0 {
		t.Errorf("reported %d findings for a handler that reads no permission:\n%v", len(got), got)
	}
}

// gin copies the parent's handler chain into a subgroup, so a group carved out
// of a guarded one is guarded. Reporting it would be a false positive, and a
// check that cries wolf is one people switch off.
func TestASubgroupInheritsTheMiddleware(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/apis/sys_user.go": apisFile,
		"app/admin/router/sys_user.go": `package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
)

func register(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	parent := v1.Group("/sys").Use(actions.PermissionAction())
	child := parent.Group("/user")
	{
		child.GET("/:id", api.Scoped)
	}
}
`,
	})
	if got := only(t, check(t, root, options{}), checkDataScopeRoute); len(got) != 0 {
		t.Errorf("reported %d findings for a subgroup of a guarded group:\n%v", len(got), got)
	}
}

// Two packages can both declare a SysUser. Only the one whose method reads the
// permission may be reported, or the check becomes a name search.
func TestAHandlerIsMatchedByPackageNotJustName(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/apis/sys_user.go": apisFile,
		"app/other/apis/sys_user.go": `package apis

import "github.com/gin-gonic/gin"

type SysUser struct{}

func (e SysUser) Scoped(c *gin.Context) {}
`,
		"app/other/router/sys_user.go": `package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/other/apis"
)

func register(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("/other")
	{
		r.GET("/:id", api.Scoped)
	}
}
`,
	})
	if got := only(t, check(t, root, options{}), checkDataScopeRoute); len(got) != 0 {
		t.Errorf("reported %d findings for a same-named handler in another package:\n%v", len(got), got)
	}
}
