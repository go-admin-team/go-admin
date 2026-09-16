package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"gorm.io/gorm"

	"go-admin/common/middleware"
)

const emptyTableNameMsg = "table name cannot be empty！"

// bodyOf covers both the success and the CustomError shape: both carry msg.
type bodyOf struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// newEngine wires one generator handler the way the router does, including the
// middleware that turns pkg.Assert's panic into a response.
//
// The generator's queries target MySQL's information_schema and cannot run on
// the sqlite connection behind them; the driver setting only has to select that
// branch, since no statement here is expected to succeed. That makes this
// serviceable for any handler in this package whose behaviour is decided before
// the query goes out -- which is what these tests are about.
func newEngine(t *testing.T, method, path string, h gin.HandlerFunc) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	previous := config.DatabaseConfig.Driver
	config.DatabaseConfig.Driver = "mysql"
	t.Cleanup(func() { config.DatabaseConfig.Driver = previous })

	r := gin.New()
	r.Use(middleware.CustomError)
	r.Handle(method, path, func(c *gin.Context) {
		c.Set("db", db)
		c.Set(pkg.LoggerKey, logger.NewHelper(logger.DefaultLogger))
		h(c)
	})
	return r
}

func newColumnListEngine(t *testing.T) *gin.Engine {
	t.Helper()
	return newEngine(t, http.MethodGet, "/db/columns/page", Gen{}.GetDBColumnList)
}

// serveJSON runs one request through the engine and decodes the envelope every
// handler here answers with. A body that will not decode fails the test rather
// than being reported as a mismatched message, which reads as the handler
// having answered something unexpected instead of not having answered at all.
func serveJSON(t *testing.T, r *gin.Engine, req *http.Request) bodyOf {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var body bodyOf
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %q: %v", w.Body.String(), err)
	}
	return body
}

func columnListMsg(t *testing.T, r *gin.Engine, query string) bodyOf {
	t.Helper()
	return serveJSON(t, r, httptest.NewRequest(http.MethodGet, "/db/columns/page"+query, nil))
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
