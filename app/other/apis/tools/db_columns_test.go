package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"

	"go-admin/common/middleware"
)

const emptyTableNameMsg = "table name cannot be empty！"

// bodyOf covers both the success and the CustomError shape: both carry msg.
type bodyOf struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// newColumnListEngine wires the handler the way the router does, including the
// middleware that turns pkg.Assert's panic into a response.
func newColumnListEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	// The query targets MySQL's information_schema; the driver setting only has
	// to select that branch, the statement itself is never expected to succeed.
	previous := config.DatabaseConfig.Driver
	config.DatabaseConfig.Driver = "mysql"
	t.Cleanup(func() { config.DatabaseConfig.Driver = previous })

	r := gin.New()
	r.Use(middleware.CustomError)
	r.GET("/db/columns/page", func(c *gin.Context) {
		c.Set("db", db)
		c.Set(pkg.LoggerKey, logger.NewHelper(logger.DefaultLogger))
		Gen{}.GetDBColumnList(c)
	})
	return r
}

func columnListMsg(t *testing.T, r *gin.Engine, query string) bodyOf {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/db/columns/page"+query, nil))

	var body bodyOf
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %q: %v", w.Body.String(), err)
	}
	return body
}

func TestGetDBColumnList_AcceptsATableName(t *testing.T) {
	body := columnListMsg(t, newColumnListEngine(t), "?tableName=sys_user")
	if body.Msg == emptyTableNameMsg {
		t.Fatalf("request carried a table name and was still rejected as empty: %+v", body)
	}
}

func TestGetDBColumnList_RejectsAMissingTableName(t *testing.T) {
	body := columnListMsg(t, newColumnListEngine(t), "")
	if body.Msg != emptyTableNameMsg {
		t.Fatalf("missing table name should be rejected, got %+v", body)
	}
}
