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
