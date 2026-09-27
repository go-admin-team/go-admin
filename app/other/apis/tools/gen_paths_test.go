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

	"go-admin/app/other/models/tools"
)

// The generator builds the paths it writes to from three fields of a table's
// configuration. These tests hold that none of them can take a write outside
// the working directory or gen.frontpath, and that a refused write leaves
// nothing behind - inside the roots or outside them.

// filesUnder lists every file below dir, or none when dir does not exist.
func filesUnder(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	return files
}

func TestNOActionsGenWritesEveryFileInsideItsRoots(t *testing.T) {
	dir := genWorkspace(t)

	var ok bool
	bodies := runGen(t, nil, func(c *gin.Context) { ok = Gen{}.NOActionsGen(c, genTable()) }, nil)
	if !ok || len(bodies) != 0 {
		t.Fatalf("ok=%v, responses %+v; want true and no error response", ok, bodies)
	}
	want := []string{
		"app/admin/models/gen_err.go",
		"app/admin/apis/gen_err.go",
		"app/admin/router/gen_err.go",
		"app/admin/service/dto/gen_err.go",
		"app/admin/service/gen_err.go",
		"ui/api/admin/gen-err.ts",
		"ui/views/admin/gen-err/index.vue",
		"ui/lang/zh-CN/gen/admin/genErr.ts",
		"ui/lang/en-US/gen/admin/genErr.ts",
	}
	for _, rel := range want {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
			t.Errorf("%s was not written: %v", rel, err)
		}
	}
	if got := len(writtenFiles(t, dir)); got != len(want) {
		t.Errorf("wrote %d files, want %d: %v", got, len(want), writtenFiles(t, dir))
	}
}

// A configuration already in sys_tables reaches NOActionsGen as it was saved,
// so this is also the case of a row saved before the save-time check existed.
func TestNOActionsGenRefusesAPathFieldThatLeavesItsRoot(t *testing.T) {
	cases := map[string]func(*tools.SysTables){
		"packageName":  func(tab *tools.SysTables) { tab.PackageName = "../../escaped" },
		"tableName":    func(tab *tools.SysTables) { tab.TBName = "../../../escaped" },
		"businessName": func(tab *tools.SysTables) { tab.BusinessName = "../../../../escaped" },
		"absolute":     func(tab *tools.SysTables) { tab.PackageName = "/escaped" },
	}
	for name, spoil := range cases {
		t.Run(name, func(t *testing.T) {
			dir := genWorkspace(t)
			tab := genTable()
			spoil(&tab)

			var ok bool
			bodies := runGen(t, nil, func(c *gin.Context) { ok = Gen{}.NOActionsGen(c, tab) }, nil)

			if ok || len(bodies) != 1 || bodies[0].Code != 500 || !strings.Contains(bodies[0].Msg, "不合法") {
				t.Fatalf("ok=%v, responses %+v; want false and one 500 naming the invalid field", ok, bodies)
			}
			if files := writtenFiles(t, dir); len(files) > 0 {
				t.Errorf("wrote %v inside the workspace", files)
			}
			// Where each of the paths above would have landed.
			for _, p := range []string{filepath.Join(dir, "..", "escaped"), filepath.Join(filepath.Dir(dir), "..", "escaped"), "/escaped"} {
				if _, err := os.Stat(p); err == nil {
					t.Errorf("%s exists after a refused generation", p)
				}
			}
		})
	}
}

// A symlink under a root that points out of it is the case a name check
// cannot see: every component of the name is legal. os.Root refuses to follow
// it. The file planted where the write would land must survive unchanged.
func TestNOActionsGenDoesNotFollowASymlinkOutOfItsRoot(t *testing.T) {
	for _, link := range []string{"app/admin", "ui/api"} {
		t.Run(link, func(t *testing.T) {
			dir := genWorkspace(t)
			outside := t.TempDir()
			planted := filepath.Join(outside, "models", "gen_err.go")
			if err := os.MkdirAll(filepath.Dir(planted), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(planted, []byte("original"), 0o644); err != nil {
				t.Fatal(err)
			}
			at := filepath.Join(dir, filepath.FromSlash(link))
			if err := os.MkdirAll(filepath.Dir(at), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, at); err != nil {
				t.Fatal(err)
			}

			var ok bool
			bodies := runGen(t, nil, func(c *gin.Context) { ok = Gen{}.NOActionsGen(c, genTable()) }, nil)

			if ok || len(bodies) != 1 || bodies[0].Code != 500 || !strings.Contains(bodies[0].Msg, "生成文件写入失败") {
				t.Fatalf("ok=%v, responses %+v; want false and one 500 for the write", ok, bodies)
			}
			if b, err := os.ReadFile(planted); err != nil || string(b) != "original" {
				t.Errorf("the file outside the root reads %q (err %v), want it untouched", b, err)
			}
			if files := filesUnder(t, outside); len(files) != 1 {
				t.Errorf("files outside the root: %v, want only the planted one", files)
			}
		})
	}
}

// putTable saves a table's configuration the way the config page does.
func putTable(t *testing.T, tab tools.SysTables) genBody {
	t.Helper()
	body, _ := json.Marshal(tab)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(body)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("db", genDB(t))
	SysTable{}.Update(c)
	var res genBody
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("decoding %q: %v", w.Body.String(), err)
	}
	return res
}

func TestSavingAConfigurationRefusesAPathField(t *testing.T) {
	db := genDB(t)
	stored := genTable()
	created, err := stored.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	edited := created
	edited.PackageName = "../escaped"
	if res := putTable(t, edited); res.Code != 500 || !strings.Contains(res.Msg, "packageName") {
		t.Fatalf("response %+v; want 500 naming packageName", res)
	}

	var row tools.SysTables
	if err := db.First(&row, created.TableId).Error; err != nil {
		t.Fatal(err)
	}
	if row.PackageName != "admin" {
		t.Errorf("package_name = %q after a refused save, want it unchanged", row.PackageName)
	}
}

// What the config page and the importer produce has to keep passing.
func TestValidateGenPathFieldsAcceptsWhatTheConfigPageAccepts(t *testing.T) {
	for _, tab := range []tools.SysTables{
		{PackageName: "admin", TBName: "sys_user", BusinessName: "sysUser"},
		{PackageName: "shop2", TBName: "order2", BusinessName: "order2"},
		{PackageName: "x", TBName: "T_Mixed_1", BusinessName: "tMixed1"},
	} {
		if err := validateGenPathFields(tab); err != nil {
			t.Errorf("%+v refused: %v", tab, err)
		}
	}
}

// The importer reads table names from the database, which accepts names no
// file can carry. A list holding one is refused whole: the valid table
// beside it is not imported either.
func TestImportRefusesATableNameThatCannotBeAFileName(t *testing.T) {
	g := newGenEnv(t)
	const good, bad = "gpa_ok", "gpa-dash"
	// dropFixture does not quote the name, and the point of bad is that it
	// needs quoting, so the rows and tables are removed here.
	drop := func() {
		g.db.Unscoped().Where("table_name IN ?", []string{good, bad}).Delete(&tools.SysTables{})
		for _, table := range []string{good, bad} {
			g.db.Exec("DROP TABLE IF EXISTS `" + table + "`")
		}
	}
	drop()
	t.Cleanup(drop)
	for _, table := range []string{good, bad} {
		if err := g.db.Exec("CREATE TABLE `" + table + "` (id int NOT NULL AUTO_INCREMENT PRIMARY KEY," + auditColumns + ")").Error; err != nil {
			t.Fatalf("creating %s: %v", table, err)
		}
	}

	if code := callHandler(t, g.db, SysTable{}.Insert, "/?tables="+good+","+bad, nil); code == http.StatusOK {
		t.Fatal("the import succeeded")
	}
	var n int64
	g.db.Model(&tools.SysTables{}).Where("table_name IN ?", []string{good, bad}).Count(&n)
	if n != 0 {
		t.Errorf("%d table(s) imported from a refused list, want 0", n)
	}
}
