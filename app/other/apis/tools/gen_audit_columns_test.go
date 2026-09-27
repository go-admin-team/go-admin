package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-admin/app/other/models/tools"
)

// The config page reads a table's columns through SysTable.Get and saves them
// back through SysTable.Update. created_at and updated_at used to be filtered
// out of what Get returns, so no generated page could list, sort or filter by
// when a row was created unless someone edited sys_columns by hand. This
// drives that page's round trip, generates the table, and runs the result.
// genEnv is a generator pointed at the MySQL test database, writing into out
// and resolving its templates from this repository.
type genEnv struct {
	db   *gorm.DB
	root string
	out  string
}

func newGenEnv(t *testing.T) genEnv {
	t.Helper()
	dsn := os.Getenv(genMySQLDSNEnv)
	if dsn == "" {
		if os.Getenv("CI") != "" {
			t.Fatalf("%s is not set while CI is: the generator must not go untested", genMySQLDSNEnv)
		}
		t.Skipf("%s is not set; the importer only reads MySQL", genMySQLDSNEnv)
	}
	gin.SetMode(gin.TestMode)
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("connecting to %s: %v", genMySQLDSNEnv, err)
	}
	var dbName string
	if err := db.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil || dbName == "" {
		t.Fatalf("reading the database name: %v", err)
	}

	out := t.TempDir()
	restore := [3]string{config.DatabaseConfig.Driver, config.GenConfig.DBName, config.GenConfig.FrontPath}
	t.Cleanup(func() {
		config.DatabaseConfig.Driver, config.GenConfig.DBName, config.GenConfig.FrontPath = restore[0], restore[1], restore[2]
	})
	config.DatabaseConfig.Driver = "mysql"
	config.GenConfig.DBName = dbName
	config.GenConfig.FrontPath = filepath.Join(out, "ui")

	if err := db.AutoMigrate(&tools.SysTables{}, &tools.SysColumns{}); err != nil {
		t.Fatalf("migrating the generator's own tables: %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "template"), filepath.Join(out, "template")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(out)
	return genEnv{db: db, root: root, out: out}
}

// importTable creates table with the given column definitions and imports it
// into the generator, returning the route parameter that names it.
func (g genEnv) importTable(t *testing.T, table, columns string) gin.Params {
	t.Helper()
	dropFixture(t, g.db, table)
	if err := g.db.Exec("CREATE TABLE " + table + " (" + columns + ")").Error; err != nil {
		t.Fatalf("creating %s: %v", table, err)
	}
	t.Cleanup(func() { dropFixture(t, g.db, table) })
	if code := callHandler(t, g.db, SysTable{}.Insert, "/?tables="+table, nil); code != http.StatusOK {
		t.Fatalf("importing %s: response code %d", table, code)
	}
	var row tools.SysTables
	if err := g.db.Where("table_name = ?", table).First(&row).Error; err != nil {
		t.Fatalf("finding %s after import: %v", table, err)
	}
	return gin.Params{{Key: "tableId", Value: itoa(row.TableId)}}
}

func (g genEnv) page(t *testing.T, table string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(g.out, "ui", "views", "admin", strings.ReplaceAll(table, "_", "-"), "index.vue"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestAuditTimestampsCanBeListedAndQueried(t *testing.T) {
	g := newGenEnv(t)
	db, root, out := g.db, g.root, g.out
	const table = "gau_event"
	tableID := g.importTable(t, table, "id int NOT NULL AUTO_INCREMENT PRIMARY KEY,"+auditColumns)

	// What the config page is shown.
	var detail struct {
		List []map[string]any `json:"list"`
		Info map[string]any   `json:"info"`
	}
	callJSON(t, db, SysTable{}.Get, http.MethodGet, "", tableID, &detail)
	shown := map[string]map[string]any{}
	for _, c := range detail.List {
		shown[c["columnName"].(string)] = c
	}
	for _, name := range []string{"created_at", "updated_at"} {
		if shown[name] == nil {
			t.Errorf("the config page is not shown %s", name)
		}
	}
	// Filled in by the framework on every write; a setting on them would
	// either do nothing or break the generated model.
	for _, name := range []string{"id", "create_by", "update_by", "deleted_at"} {
		if shown[name] != nil {
			t.Errorf("the config page is shown %s", name)
		}
	}
	if t.Failed() {
		return
	}

	// What the page sends back once someone ticks created_at's list and query
	// boxes and picks "greater than or equal".
	shown["created_at"]["isList"] = "1"
	shown["created_at"]["isQuery"] = "1"
	shown["created_at"]["queryType"] = "GTE"
	shown["updated_at"]["isList"] = "1"
	// Ticked and then unticked: the checkbox's false-value.
	shown["name"]["isQuery"] = "0"
	detail.Info["columns"] = detail.List
	body, _ := json.Marshal(detail.Info)
	if code := callJSON(t, db, SysTable{}.Update, http.MethodPut, string(body), nil, nil); code != http.StatusOK {
		t.Fatalf("saving the table config: response code %d", code)
	}

	if code := callHandler(t, db, Gen{}.GenCode, "/", tableID); code != http.StatusOK {
		t.Fatalf("generation failed with response code %d", code)
	}

	vue := g.page(t, table)
	for _, want := range []string{
		`<DateCell :value="row.createdAt" />`,
		`<DateCell :value="row.updatedAt" />`,
		`v-model="table.query.createdAt"`,
	} {
		if !strings.Contains(vue, want) {
			t.Errorf("the generated page has no %s", want)
		}
	}
	for _, unwanted := range []string{"form.model.createdAt", "form.model.updatedAt"} {
		if strings.Contains(vue, unwanted) {
			t.Errorf("the generated form edits %s, which the framework fills in", unwanted)
		}
	}

	dto, err := os.ReadFile(filepath.Join(out, "app", "admin", "service", "dto", table+".go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(dto), `form:"name"`) {
		t.Error("name, saved with isQuery \"0\", is a search field of the generated list request")
	}

	generated := backendFiles(out, root, table)
	test := filepath.Join(out, "generated_audit_test.go")
	f, err := os.Create(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := auditQueryTest.Execute(f, table); err != nil {
		t.Fatal(err)
	}
	f.Close()
	generated[filepath.Join(root, "app/admin/apis/zz_generated_audit_test.go")] = test
	overlay := filepath.Join(out, "overlay.json")
	b, _ := json.Marshal(map[string]any{"Replace": generated})
	if err := os.WriteFile(overlay, b, 0o644); err != nil {
		t.Fatal(err)
	}

	goCmd(t, root, "vet", "-overlay="+overlay, "./app/admin/...")
	ran := goCmd(t, root, "test", "-overlay="+overlay, "-count=1", "-v", "-run", "^TestGeneratedAuditQuery$", "./app/admin/apis/")
	if !strings.Contains(string(ran), "--- PASS: TestGeneratedAuditQuery ") {
		t.Error("TestGeneratedAuditQuery did not run")
	}
}

// callJSON runs one handler with a JSON body and decodes the response's data
// into data when it is not nil. It returns the code the body carries.
func callJSON(t *testing.T, db *gorm.DB, h gin.HandlerFunc, method, body string, params gin.Params, data any) int {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = params
	c.Set("db", db)
	h(c)
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("decoding the response: %v (body %q)", err, w.Body.String())
	}
	if data != nil {
		if err := json.Unmarshal(res.Data, data); err != nil {
			t.Fatalf("decoding the response data: %v (body %q)", err, w.Body.String())
		}
	}
	if res.Code != http.StatusOK {
		t.Logf("response: %s", w.Body.String())
	}
	return res.Code
}

// auditQueryTest lists the generated table filtered on created_at, sent in
// the format the generated page's date picker sends (value-format
// "YYYY-MM-DD[T]HH:mm:ssZ").
var auditQueryTest = template.Must(template.New("audit").Parse(`package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGeneratedAuditQuery(t *testing.T) {
	db, err := gorm.Open(mysql.Open(os.Getenv("GO_ADMIN_TEST_MYSQL_DSN")), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("db", db.WithContext(c)); c.Next() })
	r.GET("/x", GauEvent{}.GetPage)

	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	recent := time.Now().Add(-time.Hour)
	for name, at := range map[string]time.Time{"old": old, "recent": recent} {
		if err := db.Exec("INSERT INTO {{.}} (name, created_at, updated_at) VALUES (?, ?, ?)", name, at, at).Error; err != nil {
			t.Fatal(err)
		}
	}

	list := func(query string) []map[string]any {
		t.Helper()
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x?pageIndex=1&pageSize=10"+query, nil))
		var res struct {
			Code int
			Data struct{ List []map[string]any }
		}
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil || res.Code != 200 {
			t.Fatalf("listing %q: code %d, err %v, body %s", query, res.Code, err, w.Body.String())
		}
		return res.Data.List
	}

	all := list("")
	if len(all) != 2 {
		t.Fatalf("unfiltered list has %d rows, want 2", len(all))
	}
	if all[0]["createdAt"] == nil || all[0]["updatedAt"] == nil {
		t.Errorf("a listed row carries no createdAt/updatedAt: %v", all[0])
	}

	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02T15:04:05Z07:00")
	got := list("&createdAt=" + url.QueryEscape(since))
	if len(got) != 1 || got[0]["name"] != "recent" {
		t.Errorf("createdAt >= %s listed %v, want only the recent row", since, got)
	}
}
`))
