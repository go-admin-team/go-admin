package actions_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"go-admin/common/actions"
	"go-admin/common/dto"
	"go-admin/common/models"
)

// probeSearch is probeIndexReq with a search field, so a list request can
// be told apart from another in its response.
type probeSearch struct {
	dto.Pagination `search:"-"`
	Name           string `form:"name" search:"type:exact;column:name;table:action_probe_row"`
}

func (p *probeSearch) Generate() dto.Index        { o := *p; return &o }
func (p *probeSearch) Bind(c *gin.Context) error  { return c.ShouldBind(p) }
func (p *probeSearch) GetNeedSearch() interface{} { return *p }

// probeById names rows the way every ById DTO here does.
type probeById struct {
	dto.ObjectById
}

func (s *probeById) Generate() dto.Control                   { o := *s; return &o }
func (s *probeById) GenerateM() (models.ActiveRecord, error) { return &probeRow{}, nil }

// probeControl carries both the old actions' methods and ToModel, so one
// request type serves both generations of action.
type probeControl struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s *probeControl) Bind(c *gin.Context) error { return c.ShouldBindJSON(s) }
func (s *probeControl) Generate() dto.Control     { o := *s; return &o }
func (s *probeControl) GetId() interface{}        { return s.Id }
func (s *probeControl) GenerateM() (models.ActiveRecord, error) {
	return &probeRow{Model: models.Model{Id: s.Id}, Name: s.Name}, nil
}
func (s *probeControl) ToModel() (*probeRow, error) {
	return &probeRow{Model: models.Model{Id: s.Id}, Name: s.Name}, nil
}

// probeKey names a row by a string key, which the old actions handed to GORM
// as a SQL condition.
type probeKey struct {
	Key string `uri:"id"`
}

func (s *probeKey) Bind(c *gin.Context) error { return c.ShouldBindUri(s) }
func (s *probeKey) GetId() interface{}        { return s.Key }

func probeDB(t *testing.T, name string, l gormlogger.Interface) *gorm.DB {
	t.Helper()
	if l == nil {
		l = gormlogger.Default.LogMode(gormlogger.Silent)
	}
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{Logger: l})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&probeRow{}); err != nil {
		t.Fatal(err)
	}
	return db
}

// probeEngine serves the five routes from db, as caller 7.
func probeEngine(db *gorm.DB, register func(r gin.IRoutes)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if db != nil {
			c.Set("db", db)
		}
		c.Set(jwtauth.JwtPayloadKey, jwtauth.MapClaims{"identity": float64(7)})
		c.Next()
	})
	register(r)
	return r
}

func oldRoutes(r gin.IRoutes) {
	r.GET("/x", actions.IndexAction(&probeRow{}, &probeSearch{}, func() interface{} { l := make([]probeRow, 0); return &l }))
	r.GET("/x/:id", actions.ViewAction(&probeById{}, func() interface{} { return &probeRow{} }))
	r.POST("/x", actions.CreateAction(&probeControl{}))
	r.PUT("/x", actions.UpdateAction(&probeControl{}))
	r.DELETE("/x", actions.DeleteAction(&probeById{}))
}

func newRoutes(r gin.IRoutes) {
	r.GET("/x", actions.Index[probeRow, probeSearch]())
	r.GET("/x/:id", actions.View[probeRow, probeById]())
	r.POST("/x", actions.Create[probeRow, probeControl]())
	r.PUT("/x", actions.Update[probeRow, probeControl]())
	r.DELETE("/x", actions.Delete[probeRow, probeById]())
}

var requestIDField = regexp.MustCompile(`"requestId":"[^"]*"`)

func serve(r *gin.Engine, method, path, body string) (int, string) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w.Code, requestIDField.ReplaceAllString(w.Body.String(), `"requestId":""`)
}

// The generic actions answer every request exactly as the actions they
// replace, on the same types and the same data, until the old ones go.
func TestGenericActionsAnswerAsTheOldOnesDo(t *testing.T) {
	script := []struct{ method, path, body string }{
		{"POST", "/x", `{"name":"a"}`},
		{"POST", "/x", `{"name":"b"}`},
		{"POST", "/x", `not json`},
		{"GET", "/x?pageIndex=1&pageSize=10", ""},
		{"GET", "/x?pageIndex=1&pageSize=1&name=b", ""},
		{"GET", "/x?pageIndex=2&pageSize=1", ""},
		{"GET", "/x/1", ""},
		{"GET", "/x/99", ""},
		{"GET", "/x/abc", ""},
		{"PUT", "/x", `{"id":1,"name":"a2"}`},
		{"PUT", "/x", `{"id":99,"name":"ghost"}`},
		{"PUT", "/x", `{"name":"no-id"}`},
		{"GET", "/x/1", ""},
		{"DELETE", "/x", `{"ids":[2]}`},
		{"DELETE", "/x", `{"ids":[99]}`},
		{"GET", "/x?pageIndex=1&pageSize=10", ""},
	}
	old := probeEngine(probeDB(t, t.Name()+"-old", nil), oldRoutes)
	gen := probeEngine(probeDB(t, t.Name()+"-new", nil), newRoutes)
	for i, s := range script {
		oc, ob := serve(old, s.method, s.path, s.body)
		nc, nb := serve(gen, s.method, s.path, s.body)
		if oc != nc || ob != nb {
			t.Errorf("step %d, %s %s %s:\nold %d %s\nnew %d %s", i+1, s.method, s.path, s.body, oc, ob, nc, nb)
		}
	}
}

// Every action that reads or changes existing rows applies the caller's
// data permission: Index, View, Update and Delete.
func TestGenericActionsApplyDataPermission(t *testing.T) {
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })

	for _, req := range []struct{ method, path, body string }{
		{"GET", "/x?pageIndex=1&pageSize=10", ""},
		{"GET", "/x/1", ""},
		{"PUT", "/x", `{"id":1,"name":"a"}`},
		{"DELETE", "/x", `{"ids":[1]}`},
	} {
		cl := &capturingLogger{Interface: gormlogger.Default.LogMode(gormlogger.Silent)}
		db := probeDB(t, strings.NewReplacer("/", "_", "?", "_").Replace(t.Name()+req.method+req.path), cl)
		r := probeEngine(db, func(r gin.IRoutes) {
			self := func(c *gin.Context) {
				c.Set(actions.PermissionKey, &actions.DataPermission{DataScope: actions.DataScopeSelf, UserId: 7})
			}
			r.GET("/x", self, actions.Index[probeRow, probeSearch]())
			r.GET("/x/:id", self, actions.View[probeRow, probeById]())
			r.PUT("/x", self, actions.Update[probeRow, probeControl]())
			r.DELETE("/x", self, actions.Delete[probeRow, probeById]())
		})
		serve(r, req.method, req.path, req.body)
		if !strings.Contains(cl.all(), "action_probe_row.create_by = ") {
			t.Errorf("%s %s ran without the data-permission scope:\n%s", req.method, req.path, cl.all())
		}
	}
}

// A string key is compared as a value. Handed to Where on its own, "1=1"
// is read as a SQL condition and matches every row.
func TestGenericViewMatchesAStringKeyAsAValue(t *testing.T) {
	db := probeDB(t, t.Name(), nil)
	if err := db.Create(&probeRow{Name: "a"}).Error; err != nil {
		t.Fatal(err)
	}
	r := probeEngine(db, func(r gin.IRoutes) { r.GET("/x/:id", actions.View[probeRow, probeKey]()) })

	_, body := serve(r, "GET", "/x/1=1", "")
	var res struct{ Code int }
	_ = json.Unmarshal([]byte(body), &res)
	if res.Code != http.StatusNotFound {
		t.Errorf("GET /x/1=1 answered %s; want the row not found", body)
	}
}

// Without a database in the request the old actions wrote nothing, which a
// client reads as an empty 200.
func TestGenericActionsAnswerAMissingDatabase(t *testing.T) {
	r := probeEngine(nil, newRoutes)
	for _, req := range []struct{ method, path, body string }{
		{"GET", "/x?pageIndex=1&pageSize=10", ""},
		{"GET", "/x/1", ""},
		{"POST", "/x", `{"name":"a"}`},
		{"PUT", "/x", `{"id":1}`},
		{"DELETE", "/x", `{"ids":[1]}`},
	} {
		_, body := serve(r, req.method, req.path, req.body)
		if !strings.Contains(body, `"code":500`) || !strings.Contains(body, "数据库连接获取失败") {
			t.Errorf("%s %s answered %q; want a 500 naming the connection", req.method, req.path, body)
		}
	}
}

// Concurrent requests to one route each see only what they asked for. The
// generic actions build every value per request, so there is nothing to
// share; this holds that shape.
func TestGenericIndexKeepsConcurrentRequestsApart(t *testing.T) {
	db := probeDB(t, t.Name(), nil)
	for i := 0; i < 20; i++ {
		if err := db.Create(&probeRow{Name: fmt.Sprintf("n%d", i)}).Error; err != nil {
			t.Fatal(err)
		}
	}
	r := probeEngine(db, newRoutes)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("n%d", i)
		wg.Go(func() {
			for j := 0; j < 10; j++ {
				_, body := serve(r, "GET", "/x?pageIndex=1&pageSize=10&name="+name, "")
				var res struct {
					Data struct{ List []probeRow }
				}
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Errorf("decoding %q: %v", body, err)
					return
				}
				if len(res.Data.List) != 1 || res.Data.List[0].Name != name {
					t.Errorf("asked for %s, got %+v", name, res.Data.List)
					return
				}
			}
		})
	}
	wg.Wait()
}
