package tools

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-admin/app/other/models/tools"
)

type genBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// runGen calls h and decodes every JSON body it wrote, so a handler that
// answers twice - an error, then a success - shows up as two bodies.
func runGen(t *testing.T, db *gorm.DB, h func(*gin.Context), params gin.Params) []genBody {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = params
	if db != nil {
		c.Set("db", db)
	}
	h(c)
	var bodies []genBody
	dec := json.NewDecoder(strings.NewReader(w.Body.String()))
	for dec.More() {
		var b genBody
		if err := dec.Decode(&b); err != nil {
			t.Fatalf("decoding %q: %v", w.Body.String(), err)
		}
		bodies = append(bodies, b)
	}
	return bodies
}

// genWorkspace is a working directory laid out the way the generator expects
// to find it: the real templates under template/, output written below it.
// Any template can be swapped for one that fails when executed.
func genWorkspace(t *testing.T, broken ...string) string {
	t.Helper()
	gin.SetMode(gin.TestMode)
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	src := filepath.Join(root, "template")
	err = filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dst := filepath.Join(dir, "template", strings.TrimPrefix(p, src))
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range broken {
		// Parses, then fails on execution: SysTables has no such field.
		if err := os.WriteFile(filepath.Join(dir, "template", name), []byte("{{.NoSuchField}}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	front := config.GenConfig.FrontPath
	t.Cleanup(func() { config.GenConfig.FrontPath = front })
	config.GenConfig.FrontPath = filepath.Join(dir, "ui")
	t.Chdir(dir)
	return dir
}

func genTable() tools.SysTables {
	return tools.SysTables{
		TBName: "gen_err", ClassName: "GenErr", BusinessName: "genErr", PackageName: "admin",
		ModuleName: "gen-err", PkColumn: "id", PkGoField: "Id", PkJsonField: "id",
		Columns: []tools.SysColumns{
			{ColumnName: "id", GoField: "Id", GoType: "int", JsonField: "id", Pk: true, IsPk: "1"},
			{ColumnName: "name", GoField: "Name", GoType: "string", JsonField: "name", IsInsert: "1"},
		},
	}
}

func writtenFiles(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	for _, top := range []string{"app", "ui"} {
		_ = filepath.WalkDir(filepath.Join(dir, top), func(p string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				files = append(files, p)
			}
			return nil
		})
	}
	return files
}

// A template that fails on execution used to leave its file empty and every
// other file written, with the request reporting success.
func TestNOActionsGenWritesNothingWhenATemplateFails(t *testing.T) {
	dir := genWorkspace(t, "v4/no_actions/service.go.template")

	var ok bool
	bodies := runGen(t, nil, func(c *gin.Context) { ok = Gen{}.NOActionsGen(c, genTable()) }, nil)

	if ok || len(bodies) != 1 || bodies[0].Code != 500 {
		t.Fatalf("ok=%v, responses %+v; want false and one 500", ok, bodies)
	}
	if !strings.Contains(bodies[0].Msg, "service.go.template") {
		t.Errorf("message %q does not name the template that failed", bodies[0].Msg)
	}
	if files := writtenFiles(t, dir); len(files) > 0 {
		t.Errorf("wrote %d file(s) although a template failed: %v", len(files), files)
	}
}

// pkg.FileCreate, which this used to write through, ends the process with
// log.Fatalln when it cannot create the file. Reaching the assertions at all
// is half of this test.
func TestNOActionsGenReportsAFileItCannotWrite(t *testing.T) {
	dir := genWorkspace(t)
	blocker := filepath.Join(dir, "app", "admin", "models")
	if err := os.MkdirAll(filepath.Dir(blocker), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blocker, nil, 0o644); err != nil { // a file where a directory must go
		t.Fatal(err)
	}

	var ok bool
	bodies := runGen(t, nil, func(c *gin.Context) { ok = Gen{}.NOActionsGen(c, genTable()) }, nil)

	if ok || len(bodies) != 1 || bodies[0].Code != 500 {
		t.Fatalf("ok=%v, responses %+v; want false and one 500", ok, bodies)
	}
	// Not a template failing first for some other reason.
	if !strings.Contains(bodies[0].Msg, "生成文件写入失败") {
		t.Errorf("message %q is not the write failure", bodies[0].Msg)
	}
}

func genDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(new(tools.SysTables), new(tools.SysColumns)); err != nil {
		t.Fatal(err)
	}
	return db
}

// GenCode answered the request itself after NOActionsGen had already answered
// it with an error, so the client received the error followed by
// "Code generated successfully".
func TestGenCodeAnswersOnceWhenGenerationFails(t *testing.T) {
	genWorkspace(t, "v4/no_actions/service.go.template")
	db := genDB(t)
	table := genTable()
	created, err := table.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	bodies := runGen(t, db, Gen{}.GenCode, gin.Params{{Key: "tableId", Value: itoa(created.TableId)}})

	if len(bodies) != 1 || bodies[0].Code != 500 {
		t.Errorf("responses %+v; want exactly one, a 500", bodies)
	}
}
