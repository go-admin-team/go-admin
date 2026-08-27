package version

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The repository carries two ModelTime types. The one in
// cmd/migrate/migration/models still has a nullable gorm.DeletedAt and is what
// builds the tables; the one in common/models is the millisecond marker and is
// what queries them. A table whose runtime model embeds the second but which no
// migration converts is queried with deleted_at = 0 against a datetime column,
// and every row is invisible - silently, and only in production.
//
// sys_columns and sys_tables were in exactly that state: the code generator
// listed no tables at all.
func TestEveryRuntimeSoftDeleteTableIsConverted(t *testing.T) {
	converted := map[string]bool{}
	for _, name := range append(append([]string{}, softDeleteTables...), generatorTables...) {
		converted[name] = true
	}

	// tb_demo is built with the marker-less model and has no runtime model at
	// all, so nothing ever queries it with deleted_at = 0.
	const noRuntimeModel = "tb_demo"

	for table, file := range runtimeSoftDeleteTables(t) {
		if table == noRuntimeModel {
			continue
		}
		if !converted[table] {
			t.Errorf("%s (%s) embeds common.ModelTime but no migration converts its deleted_at;\n"+
				"    it will be queried with deleted_at = 0 against a nullable datetime and return nothing",
				table, file)
		}
	}
}

// runtimeSoftDeleteTables maps table name to the file declaring it, for every
// model under app/ that embeds the marker-carrying ModelTime.
func runtimeSoftDeleteTables(t *testing.T) map[string]string {
	t.Helper()

	root := repoRoot(t)
	found := map[string]string{}

	err := filepath.Walk(filepath.Join(root, "app"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil // not this test's business
		}
		// Only files importing the runtime models package can embed its ModelTime.
		if !importsRuntimeModels(f) {
			return nil
		}
		for name, table := range tablesWithModelTime(f) {
			_ = name
			found[table] = strings.TrimPrefix(path, root+"/")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk app/: %v", err)
	}
	if len(found) == 0 {
		t.Fatal("found no runtime models at all; the scan is broken, not the code")
	}
	return found
}

func importsRuntimeModels(f *ast.File) bool {
	for _, imp := range f.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err == nil && p == "go-admin/common/models" {
			return true
		}
	}
	return false
}

// tablesWithModelTime returns struct name -> table name for structs that embed
// ModelTime and declare a TableName.
func tablesWithModelTime(f *ast.File) map[string]string {
	embeds := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range st.Fields.List {
			if len(field.Names) != 0 {
				continue // named field, not an embed
			}
			if sel, ok := field.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "ModelTime" {
				embeds[ts.Name.Name] = true
			}
		}
		return true
	})

	out := map[string]string{}
	for name := range embeds {
		if table := tableNameOf(f, name); table != "" {
			out[name] = table
		}
	}
	return out
}

// tableNameOf finds the string returned by func (T) TableName() string.
func tableNameOf(f *ast.File, structName string) string {
	var table string
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "TableName" || fn.Recv == nil || len(fn.Recv.List) != 1 {
			return true
		}
		if receiverName(fn.Recv.List[0].Type) != structName {
			return true
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if ok && lit.Kind == token.STRING {
				if s, err := strconv.Unquote(lit.Value); err == nil && table == "" {
					table = s
				}
			}
			return true
		})
		return true
	})
	return table
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	}
	return ""
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("go.mod not found above the test directory")
	return ""
}
