package actions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"

	"go-admin/common/actions"
)

// The detailed data-permission regression suite (claims parsing, the
// GetOrm-unavailable abort, the SQL each data scope produces) now lives in
// go-admin-core's sdk/contract/actions, alongside the logic itself (PRD 006
// F3). What is left to test here is the shim's own wiring: that this
// package's exported names still round-trip through the same *gin.Context
// key core's PermissionAction and GetPermissionFromContext use.
//
// This file lives in package actions_test, an external test, deliberately:
// it exercises PermissionAction and GetPermissionFromContext exactly as an
// app/admin Service does, through this package's public API only, not
// through anything internal a forward could paper over.

// TestPermissionKeyMatchesWhatPermissionActionSets guards PRD 006's hard
// constraint 4: PermissionKey must be declared as
// `const PermissionKey = contractactions.PermissionKey`, a direct
// reference, never a restated literal (see type.go). PermissionAction is
// core's middleware and always writes under core's own key. This test reads
// the value back with actions.PermissionKey exactly as code outside
// GetPermissionFromContext would - c.Get(actions.PermissionKey) is a real,
// if uncommon, way to read the value go-admin has always allowed, and it is
// the one call site where an independently declared PermissionKey would
// stop working without GetPermissionFromContext's own forward hiding it.
//
// If PermissionKey were ever re-declared as an independent literal in this
// package, a later edit to core's copy would make this test fail without a
// single byte of this package having changed - which is the silent-failure
// mode hard constraint 4 exists to rule out (evaluation S2).
func TestPermissionKeyMatchesWhatPermissionActionSets(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set(jwt.JwtPayloadKey, jwt.MapClaims{
		"identity":  float64(7),
		"roleid":    float64(3),
		"deptid":    float64(5),
		"datascope": actions.DataScopeDeptTree,
	})

	actions.PermissionAction()(c)

	value, ok := c.Get(actions.PermissionKey)
	if !ok {
		t.Fatal("PermissionAction did not set the key actions.PermissionKey names; the two have diverged")
	}
	p, ok := value.(*actions.DataPermission)
	if !ok || p.DataScope != actions.DataScopeDeptTree || p.DeptId != 5 {
		t.Fatalf("value under actions.PermissionKey = %#v, want a DataPermission carrying the token's scope", value)
	}
}

// TestGetPermissionFromContextRoundTrips is the same guard from the other
// exported entry point: GetPermissionFromContext must read back exactly
// what PermissionAction wrote, both reached through this package's own
// forwards rather than core's directly.
func TestGetPermissionFromContextRoundTrips(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set(jwt.JwtPayloadKey, jwt.MapClaims{
		"identity":  float64(7),
		"roleid":    float64(3),
		"deptid":    float64(5),
		"datascope": actions.DataScopeSelf,
	})

	actions.PermissionAction()(c)

	p := actions.GetPermissionFromContext(c)
	if p.DataScope != actions.DataScopeSelf || p.UserId != 7 {
		t.Fatalf("GetPermissionFromContext() = %+v, want DataScope=%q UserId=7", p, actions.DataScopeSelf)
	}
}

func TestIsValidDataScope(t *testing.T) {
	for _, s := range []string{actions.DataScopeAll, actions.DataScopeCustom, actions.DataScopeDept, actions.DataScopeDeptTree, actions.DataScopeSelf} {
		if !actions.IsValidDataScope(s) {
			t.Errorf("IsValidDataScope(%q) = false, want true", s)
		}
	}
	if actions.IsValidDataScope("6") {
		t.Error(`IsValidDataScope("6") = true, want false`)
	}
}
