package actions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
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

// PRD 006 F14/H1. A token without deptid/datascope forces the fallback
// query, which needs pkg.GetOrm(c) - and no "db" key is set in this
// context, so GetOrm fails exactly as it would if a tenant's database were
// unreachable. Before the fix, that error was logged and the handler ran
// anyway with no data permission filter at all.
func TestPermissionActionAbortsWhenDBIsUnavailable(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handlerReached := false
	r.Use(func(c *gin.Context) {
		c.Set(jwt.JwtPayloadKey, jwt.MapClaims{"identity": float64(7)})
	})
	r.Use(PermissionAction())
	r.GET("/", func(c *gin.Context) {
		handlerReached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if handlerReached {
		t.Fatal("the business handler ran with no database and no data permission filter set")
	}
}

// PRD 006 F14/H2 and H3. Table-driven over gorm DryRun so the exact SQL
// Permission produces for each scope is pinned down, not just "some WHERE
// clause got added".
func TestPermissionScopes(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	const noRows = "SELECT * FROM `t` WHERE 1 = 0"

	cases := []struct {
		name string
		p    *DataPermission
		want string
		vars []interface{}
	}{
		{
			name: "all",
			p:    &DataPermission{DataScope: DataScopeAll},
			want: "SELECT * FROM `t`",
		},
		{
			name: "custom",
			p:    &DataPermission{DataScope: DataScopeCustom, RoleId: 3},
			want: "SELECT * FROM `t` WHERE t.create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)",
			vars: []interface{}{3},
		},
		{
			name: "dept",
			p:    &DataPermission{DataScope: DataScopeDept, DeptId: 5},
			want: "SELECT * FROM `t` WHERE t.create_by in (SELECT user_id from sys_user where dept_id = ? )",
			vars: []interface{}{5},
		},
		{
			name: "dept-tree",
			p:    &DataPermission{DataScope: DataScopeDeptTree, DeptId: 5},
			want: "SELECT * FROM `t` WHERE t.create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))",
			vars: []interface{}{"%/5/%"},
		},
		{
			name: "self",
			p:    &DataPermission{DataScope: DataScopeSelf, UserId: 7},
			want: "SELECT * FROM `t` WHERE t.create_by = ?",
			vars: []interface{}{7},
		},
		// H2: an unrecognized scope must not read like "all data" any more.
		{name: "unrecognized value", p: &DataPermission{DataScope: "6"}, want: noRows},
		// H2/H1: the zero-value DataPermission is what getPermissionFromContext
		// and the two "give up and continue" branches in PermissionAction hand
		// out when nothing else is available.
		{name: "zero value (no scope at all)", p: &DataPermission{}, want: noRows},
		// H3: dept_path always starts with "/0/" (sys_dept.go), so DeptId 0
		// must not be allowed to build a pattern that matches every row.
		{name: "dept with DeptId 0", p: &DataPermission{DataScope: DataScopeDept, DeptId: 0}, want: noRows},
		{name: "dept-tree with DeptId 0", p: &DataPermission{DataScope: DataScopeDeptTree, DeptId: 0}, want: noRows},
		{name: "dept-tree with negative DeptId", p: &DataPermission{DataScope: DataScopeDeptTree, DeptId: -1}, want: noRows},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := db.Session(&gorm.Session{DryRun: true}).
				Table("t").
				Scopes(Permission("t", tc.p)).
				Find(&[]map[string]interface{}{}).
				Statement

			if stmt.SQL.String() != tc.want {
				t.Errorf("SQL = %q, want %q", stmt.SQL.String(), tc.want)
			}
			if tc.vars == nil {
				if len(stmt.Vars) != 0 {
					t.Errorf("vars = %v, want none", stmt.Vars)
				}
				return
			}
			if len(stmt.Vars) != len(tc.vars) {
				t.Fatalf("vars = %v, want %v", stmt.Vars, tc.vars)
			}
			for i := range tc.vars {
				if stmt.Vars[i] != tc.vars[i] {
					t.Errorf("vars[%d] = %v, want %v", i, stmt.Vars[i], tc.vars[i])
				}
			}
		})
	}
}
