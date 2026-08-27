package version

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Fields tagged gorm:"size:4" become a tinyint on MySQL, which holds -128..127.
// sqlite ignores the width, so a value that overflows passes every local test
// and fails on a real install - and because the migration is not transactional,
// it fails partway, leaving later migrations unapplied.
//
// That is what happened: a seeded menu with Sort: 900 stopped the run at
// 1786700001000, so the soft-delete conversion never ran, deleted_at stayed
// NULL, and nobody could log in.
var narrowColumns = map[string]struct{ min, max int64 }{
	"Sort":   {math.MinInt8, math.MaxInt8},
	"Status": {math.MinInt8, math.MaxInt8},
}

func TestSeededValuesFitTheirColumns(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	checked := 0
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			kv, ok := n.(*ast.KeyValueExpr)
			if !ok {
				return true
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				return true
			}
			limits, watched := narrowColumns[key.Name]
			if !watched {
				return true
			}
			lit, ok := kv.Value.(*ast.BasicLit)
			if !ok || lit.Kind != token.INT {
				return true
			}
			v, err := strconv.ParseInt(lit.Value, 10, 64)
			if err != nil {
				return true
			}
			checked++
			if v < limits.min || v > limits.max {
				t.Errorf("%s:%d: %s: %d does not fit a tinyint (%d..%d);\n"+
					"    MySQL rejects it with Error 1264 and the migration stops there",
					filepath.Base(path), fset.Position(lit.Pos()).Line, key.Name, v, limits.min, limits.max)
			}
			return true
		})
	}

	if checked == 0 {
		t.Fatal("no seeded values were examined; the scan is broken, not the code")
	}
}
