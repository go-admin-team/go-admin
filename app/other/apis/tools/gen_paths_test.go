package tools

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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
