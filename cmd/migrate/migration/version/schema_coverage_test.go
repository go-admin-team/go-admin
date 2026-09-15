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

	"go-admin/cmd/migrate/migration"
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
	return importsPackage(f, "go-admin/common/models")
}

func importsPackage(f *ast.File, pkg string) bool {
	for _, imp := range f.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err == nil && p == pkg {
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

// softDeleteConversion is the version at which sys_api, sys_menu and the rest
// stop storing deleted_at as a nullable timestamp and start storing the NOT
// NULL millisecond marker.
const softDeleteConversion = 1786700003000

// versionPrefixLen is the width migration.GetFilename slices off a filename.
const versionPrefixLen = 13

// Migrations ordered after the conversion must not seed rows through
// cmd/migrate/migration/models.
//
// That package's ModelTime still declares a nullable gorm.DeletedAt, which is
// correct for the migrations that predate the conversion - it is the shape the
// column had when they ran. Reusing it afterwards writes NULL into a NOT NULL
// column and the migration fails on its first insert:
//
//	NOT NULL constraint failed: sys_api.deleted_at
//
// A fresh database never catches this, because every migration using that
// package today is ordered before the conversion and so runs while the column
// is still nullable. Only a migration added afterwards hits it, which in
// practice means the next person adding a business module - the reference
// they copy, 1786700001000_demo_menu.go, is itself one of the safe ones.
//
// Reads through that package are worse than writes, which is why the whole
// import is banned rather than just the inserts. gorm scopes a nullable
// DeletedAt as "WHERE deleted_at IS NULL", and after the conversion live rows
// hold 0, so the row is simply not there:
//
//	frozen  SysRole -> record not found
//	runtime SysRole -> roleId=1
//
// 1786700001000_demo_menu.go looks the admin role up that way and treats
// ErrRecordNotFound as "roles are not seeded yet, skip authorisation". A
// post-conversion copy that switched its inserts to the runtime models but
// kept this lookup would seed the menu, grant nothing, and still record the
// migration as applied - the menu appears, its buttons do nothing, and no
// error is reported anywhere.
func TestPostConversionMigrationsAvoidFrozenSeedModels(t *testing.T) {
	const frozenModels = "go-admin/cmd/migrate/migration/models"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	checked := 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		// GetFilename is what every migration uses to derive its own version,
		// so the two stay in step if the filename convention ever changes.
		if len(name) < versionPrefixLen {
			continue
		}
		version, err := strconv.ParseInt(migration.GetFilename(name), 10, 64)
		if err != nil || version <= softDeleteConversion {
			continue // not a versioned migration, or one that predates the change
		}
		checked++

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		if importsPackage(f, frozenModels) {
			t.Errorf("%s is ordered after the soft-delete conversion but seeds through %s;\n"+
				"    that package writes a nullable deleted_at and will fail with\n"+
				"    \"NOT NULL constraint failed\" on its first insert.\n"+
				"    Use the runtime models under app/ instead - they carry the marker.",
				name, frozenModels)
		}
	}

	if checked == 0 {
		t.Fatal("no post-conversion migrations found; the scan is broken, not the code")
	}
}
