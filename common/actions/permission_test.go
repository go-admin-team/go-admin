package actions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// No database is placed in the context on purpose. The middleware needs one
// only to run the sys_user join, so reaching the handler proves it did not.
func runPermission(t *testing.T, claims jwt.MapClaims) (*DataPermission, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if claims != nil {
		c.Set(jwt.JwtPayloadKey, claims)
	}

	PermissionAction()(c)

	value, exists := c.Get(PermissionKey)
	if !exists {
		return nil, false
	}
	p, _ := value.(*DataPermission)
	return p, true
}

// Permission() returns the query untouched when data permission is off, so the
// lookup feeding it has nothing to feed. It used to run regardless: a sys_user
// join on every list, detail, update and delete, discarded immediately.
func TestNoLookupWhenDataPermissionIsOff(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = false
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	if _, ok := runPermission(t, jwt.MapClaims{"identity": float64(7)}); !ok {
		t.Fatal("the request needed a database even though data permission is off")
	}
}

func TestScopeComesFromTheTokenWhenItCarriesOne(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	p, ok := runPermission(t, jwt.MapClaims{
		"identity":  float64(7),
		"roleid":    float64(3),
		"deptid":    float64(5),
		"datascope": "4",
	})
	if !ok {
		t.Fatal("the token carried the scope and a database was still needed")
	}
	if p.DataScope != "4" || p.UserId != 7 || p.DeptId != 5 || p.RoleId != 3 {
		t.Fatalf("scope read as %+v", p)
	}
}

// A token minted before deptid was carried is still valid until it expires, and
// has to keep working - by falling back to the query, which needs a database.
func TestATokenWithoutDeptIdFallsBackToTheQuery(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	if _, ok := runPermission(t, jwt.MapClaims{
		"identity":  float64(7),
		"roleid":    float64(3),
		"datascope": "4",
	}); ok {
		t.Fatal("an old token was served from claims it does not have")
	}
}
