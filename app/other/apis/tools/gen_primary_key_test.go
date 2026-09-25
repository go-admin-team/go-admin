package tools

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
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

// The importer reads tables from information_schema, which only MySQL is
// supported for, so this runs only when a MySQL DSN is set - and must run in
// CI, where one is. A suite that skipped there would report the generator
// fine for every primary key it has never been shown.
const genMySQLDSNEnv = "GO_ADMIN_TEST_MYSQL_DSN"

// auditColumns are the columns every generated model maps through
// models.ModelTime and models.ControlBy, typed the way those structs store
// them. A fixture without them would describe a table the generated model
// cannot read or write.
const auditColumns = `name varchar(64) NOT NULL COMMENT 'name',
	create_by bigint NULL,
	update_by bigint NULL,
	created_at datetime(3) NULL,
	updated_at datetime(3) NULL,
	deleted_at bigint NOT NULL DEFAULT 0`

// pkShape is one table the importer can be pointed at. The shapes are the
// primary keys a real table has, not the one the templates were written
// for.
type pkShape struct {
	Table    string
	Class    string
	Key      string // key column definition, before the audit columns
	Tail     string // table constraints, after them
	Refused  bool   // generation must refuse the table outright
	PkColumn string
	// Assigned is true when the database assigns the key; otherwise Create
	// carries it and ID is the value it carries.
	Assigned bool
	Create   string
	ID       string
	IDJSON   string
	Missing  string
}

var pkShapes = []pkShape{{
	Table: "gpk_serial", Class: "GpkSerial", PkColumn: "id",
	Key:      "id int NOT NULL AUTO_INCREMENT PRIMARY KEY,",
	Assigned: true, Create: `{"name":"alpha"}`, Missing: "987654",
}, {
	Table: "gpk_order", Class: "GpkOrder", PkColumn: "order_no",
	Key:      "order_no int NOT NULL AUTO_INCREMENT PRIMARY KEY,",
	Assigned: true, Create: `{"name":"alpha"}`, Missing: "987654",
}, {
	Table: "gpk_code", Class: "GpkCode", PkColumn: "code",
	Key:    "code varchar(32) NOT NULL PRIMARY KEY,",
	Create: `{"code":"c-1","name":"alpha"}`, ID: "c-1", IDJSON: `"c-1"`, Missing: "no-such-code",
}, {
	Table: "gpk_nokey", Key: "ref int NOT NULL,", Refused: true,
}, {
	Table: "gpk_pair", Key: "a int NOT NULL, b int NOT NULL,", Tail: ", PRIMARY KEY (a, b)", Refused: true,
}}

func TestGeneratorHandlesEveryPrimaryKeyShape(t *testing.T) {
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
	for _, s := range pkShapes {
		dropFixture(t, db, s.Table)
		ddl := "CREATE TABLE " + s.Table + " (" + s.Key + auditColumns + s.Tail + ")"
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("creating %s: %v", s.Table, err)
		}
		t.Cleanup(func() { dropFixture(t, db, s.Table) })
	}

	// The generator resolves its templates, and writes its output, relative to
	// the working directory, exactly as it does in a running server.
	if err := os.Symlink(filepath.Join(root, "template"), filepath.Join(out, "template")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(out)

	generated := map[string]string{}
	var runnable []pkShape
	genFailed := false
	for _, s := range pkShapes {
		if code := callHandler(t, db, SysTable{}.Insert, "/?tables="+s.Table, nil); code != http.StatusOK {
			t.Fatalf("importing %s: response code %d", s.Table, code)
		}
		var row tools.SysTables
		if err := db.Where("table_name = ?", s.Table).First(&row).Error; err != nil {
			t.Fatalf("finding %s after import: %v", s.Table, err)
		}
		code := callHandler(t, db, Gen{}.GenCode, "/", gin.Params{{Key: "tableId", Value: itoa(row.TableId)}})
		files := backendFiles(out, root, s.Table)

		if s.Refused {
			if code == http.StatusOK {
				t.Errorf("%s: generation succeeded; a table without exactly one primary key column must be refused", s.Table)
			}
			for _, generatedFile := range files {
				if _, err := os.Stat(generatedFile); err == nil {
					t.Errorf("%s: refused, but %s was written anyway", s.Table, generatedFile)
				}
			}
			continue
		}
		if code != http.StatusOK {
			t.Errorf("%s: generation failed with response code %d", s.Table, code)
			genFailed = true
			continue
		}
		for dst, src := range files {
			generated[dst] = src
		}
		runnable = append(runnable, s)
	}
	if genFailed {
		return
	}

	// Compile the generated files in place, inside this module, through an
	// overlay: nothing is written into the working tree, and the generated
	// code is checked against the real packages it will live in.
	crud := filepath.Join(out, "generated_crud_test.go")
	f, err := os.Create(crud)
	if err != nil {
		t.Fatal(err)
	}
	if err := crudTest.Execute(f, runnable); err != nil {
		t.Fatal(err)
	}
	f.Close()
	generated[filepath.Join(root, "app/admin/apis/zz_generated_crud_test.go")] = crud

	overlay := filepath.Join(out, "overlay.json")
	b, _ := json.Marshal(map[string]any{"Replace": generated})
	if err := os.WriteFile(overlay, b, 0o644); err != nil {
		t.Fatal(err)
	}

	goCmd(t, root, "vet", "-overlay="+overlay, "./app/admin/...")
	ran := goCmd(t, root, "test", "-overlay="+overlay, "-count=1", "-v", "-run", "^TestGeneratedCRUD_", "./app/admin/apis/")
	// A -run pattern that matches nothing also exits 0. Each table has to
	// have been driven, not merely compiled.
	for _, s := range runnable {
		if !strings.Contains(string(ran), "--- PASS: TestGeneratedCRUD_"+s.Class+" ") {
			t.Errorf("%s: TestGeneratedCRUD_%s did not run", s.Table, s.Class)
		}
	}
}

func dropFixture(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	var ids []int
	db.Model(&tools.SysTables{}).Where("table_name = ?", table).Pluck("table_id", &ids)
	if len(ids) > 0 {
		db.Where("table_id IN ?", ids).Delete(&tools.SysColumns{})
		db.Where("table_id IN ?", ids).Delete(&tools.SysTables{})
	}
	if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
		t.Fatalf("dropping %s: %v", table, err)
	}
}

// callHandler runs one handler the way a request would reach it and returns
// the code its first response body carries. GenCode can write a second body
// after an error body, so only the first is the verdict.
func callHandler(t *testing.T, db *gorm.DB, h gin.HandlerFunc, target string, params gin.Params) int {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, target, nil)
	c.Params = params
	c.Set("db", db)
	h(c)
	var body struct {
		Code int `json:"code"`
	}
	if err := json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&body); err != nil {
		t.Fatalf("decoding the response to %s: %v (body %q)", target, err, w.Body.String())
	}
	return body.Code
}

// backendFiles maps where each generated backend file belongs in this module
// to where the generator wrote it.
func backendFiles(out, root, table string) map[string]string {
	m := map[string]string{}
	for _, dir := range []string{"models", "apis", "router", "service", "service/dto"} {
		rel := filepath.Join("app/admin", dir, table+".go")
		m[filepath.Join(root, rel)] = filepath.Join(out, rel)
	}
	return m
}

func goCmd(t *testing.T, dir string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		t.Fatalf("go %s failed:\n%s", strings.Join(args, " "), out)
	}
	if err != nil {
		t.Fatalf("running go %s: %v", args[0], err)
	}
	t.Logf("go %s:\n%s", args[0], out)
	return out
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

// crudTest drives each generated table through its generated handlers, over
// HTTP, against the fixture table. The routes use the same "/:id" pattern the
// generated routers register.
var crudTest = template.Must(template.New("crud").Parse(`package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type generatedResp struct {
	Code int             ` + "`json:\"code\"`" + `
	Data json.RawMessage ` + "`json:\"data\"`" + `
}

func generatedEngine(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(mysql.Open(os.Getenv("GO_ADMIN_TEST_MYSQL_DSN")), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("db", db.WithContext(c)); c.Next() })
	return r, db
}

func generatedDo(r *gin.Engine, method, path, body string) generatedResp {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var res generatedResp
	_ = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&res)
	return res
}

func generatedLive(db *gorm.DB, table string) int64 {
	var n int64
	db.Table(table).Where("deleted_at = 0").Count(&n)
	return n
}
{{range .}}
func TestGeneratedCRUD_{{.Class}}(t *testing.T) {
	r, db := generatedEngine(t)
	a := {{.Class}}{}
	r.POST("/x", a.Insert)
	r.GET("/x/:id", a.Get)
	r.PUT("/x/:id", a.Update)
	r.DELETE("/x", a.Delete)

	if res := generatedDo(r, http.MethodPost, "/x", ` + "`{{.Create}}`" + `); res.Code != 200 {
		t.Fatalf("insert: code %d", res.Code)
	}
{{- if .Assigned}}
	var n int
	db.Table("{{.Table}}").Select("MAX({{.PkColumn}})").Scan(&n)
	id, idJSON := itoaGenerated(n), itoaGenerated(n)
{{- else}}
	id, idJSON := "{{.ID}}", ` + "`{{.IDJSON}}`" + `
{{- end}}

	t.Run("get by the path parameter", func(t *testing.T) {
		res := generatedDo(r, http.MethodGet, "/x/"+url.PathEscape(id), "")
		var row map[string]any
		_ = json.Unmarshal(res.Data, &row)
		if res.Code != 200 || row["name"] != "alpha" {
			t.Errorf("code %d, data %s", res.Code, res.Data)
		}
	})
	t.Run("update by the path parameter updates, not inserts", func(t *testing.T) {
		live := generatedLive(db, "{{.Table}}")
		res := generatedDo(r, http.MethodPut, "/x/"+url.PathEscape(id), ` + "`{\"name\":\"beta\"}`" + `)
		var name string
		db.Table("{{.Table}}").Where("{{.PkColumn}} = ?", id).Select("name").Scan(&name)
		if res.Code != 200 || name != "beta" || generatedLive(db, "{{.Table}}") != live {
			t.Errorf("code %d, name %q, live rows %d want %d", res.Code, name, generatedLive(db, "{{.Table}}"), live)
		}
	})
	t.Run("a missing key is not found", func(t *testing.T) {
		if res := generatedDo(r, http.MethodGet, "/x/{{.Missing}}", ""); res.Code == 200 {
			t.Errorf("found: %s", res.Data)
		}
	})
	t.Run("updating a missing key does not create a row", func(t *testing.T) {
		live := generatedLive(db, "{{.Table}}")
		res := generatedDo(r, http.MethodPut, "/x/{{.Missing}}", ` + "`{\"name\":\"ghost\"}`" + `)
		if got := generatedLive(db, "{{.Table}}"); got != live || res.Code == 200 {
			t.Errorf("code %d, live rows %d want %d", res.Code, got, live)
		}
	})
	// "1=1" rather than a quote-breaking payload: a key handed to GORM as a
	// bare string becomes the WHERE clause itself, and this one is valid SQL
	// that matches every row, so the failure it exposes is a successful read.
	t.Run("a key carrying SQL is a value, not a condition", func(t *testing.T) {
		if res := generatedDo(r, http.MethodGet, "/x/"+url.PathEscape("1=1"), ""); res.Code == 200 {
			t.Errorf("found: %s", res.Data)
		}
	})
	t.Run("delete by key", func(t *testing.T) {
		live := generatedLive(db, "{{.Table}}")
		res := generatedDo(r, http.MethodDelete, "/x", ` + "`{\"ids\":[`+idJSON+`]}`" + `)
		if res.Code != 200 || generatedLive(db, "{{.Table}}") != live-1 {
			t.Errorf("code %d, live rows %d want %d", res.Code, generatedLive(db, "{{.Table}}"), live-1)
		}
	})
}
{{end}}
func itoaGenerated(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}
`))
