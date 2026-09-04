package actions_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"

	"go-admin/common/actions"
	"go-admin/common/dto"
	"go-admin/common/models"
)

// capturingLogger records every SQL statement GORM actually executes, so a
// test can inspect it the way inspecting a *gorm.DB's own Statement cannot:
// IndexAction builds and executes its query in one unbroken chain
// (db.Model(...).Scopes(...).Find(...)...Count(...)) and never hands the
// built statement back to its caller.
type capturingLogger struct {
	gormlogger.Interface
	mu    sync.Mutex
	stmts []string
}

func (l *capturingLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	l.mu.Lock()
	l.stmts = append(l.stmts, sql)
	l.mu.Unlock()
}

func (l *capturingLogger) all() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.stmts, "\n")
}

// probeRow is a minimal model satisfying models.ActiveRecord through the
// same embeds a real app/admin model uses, so IndexAction sees exactly the
// shape it is written against.
type probeRow struct {
	models.Model
	models.ControlBy
	Name string
}

func (probeRow) TableName() string                { return "action_probe_row" }
func (e *probeRow) Generate() models.ActiveRecord { o := *e; return &o }
func (e *probeRow) GetId() interface{}            { return e.Id }

// probeIndexReq is a minimal dto.Index: no search tags, page defaults.
type probeIndexReq struct {
	dto.Pagination `search:"-"`
}

func (p *probeIndexReq) Generate() dto.Index        { return p }
func (p *probeIndexReq) Bind(*gin.Context) error    { return nil }
func (p *probeIndexReq) GetNeedSearch() interface{} { return *p }

type pageEnvelope struct {
	Code int32 `json:"code"`
}

// TestIndexActionAppliesDataPermission is an end-to-end guard core's own
// test suite cannot provide. The five generic CRUD actions in this package
// (create/delete/index/update/view.go) were not lowered to core (PRD 006
// F3) - they still call actions.Permission directly, in this repository, on
// a code path core knows nothing about. core's tests pin down what
// Permission does for a given scope; nothing pinned down whether this
// package's own Actions still remember to call it at all. This runs
// IndexAction exactly as a real request would, against a real in-memory
// database, and inspects the SQL GORM actually executed - not just that
// the handler returned success, which it would just as happily do with no
// filter applied at all.
func TestIndexActionAppliesDataPermission(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	cl := &capturingLogger{Interface: gormlogger.Default.LogMode(gormlogger.Silent)}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: cl})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&probeRow{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("db", db)
	c.Set(actions.PermissionKey, &actions.DataPermission{DataScope: actions.DataScopeSelf, UserId: 7})

	actions.IndexAction(&probeRow{}, &probeIndexReq{}, func() interface{} { return &[]probeRow{} })(c)

	var body pageEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body %q: %v", w.Body.String(), err)
	}
	if body.Code != http.StatusOK {
		t.Fatalf("response code = %d, want %d; body=%s", body.Code, http.StatusOK, w.Body.String())
	}

	sql := cl.all()
	const wantFragment = "action_probe_row.create_by = "
	if !strings.Contains(sql, wantFragment) {
		t.Fatalf("IndexAction did not apply the data-permission scope to its query; want SQL containing %q, got:\n%s", wantFragment, sql)
	}
}
