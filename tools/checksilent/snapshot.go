package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// skippedDirs are never walked. Generated and vendored trees would produce
// findings nobody can act on, and node_modules would make the walk the slowest
// part of CI.
var skippedDirs = map[string]bool{
	".git":         true,
	".idea":        true,
	"node_modules": true,
	"vendor":       true,
	"testdata":     true,
	"dist":         true,
	"temp":         true,
}

// sourceFile is one parsed Go file plus the mapping from the local names it
// uses for packages to the paths those names stand for.
//
// Resolving through the file's own import block, rather than matching on the
// identifier, is what keeps service.SysMenu apart from models.SysMenu.
type sourceFile struct {
	// Path is relative to the scanned root and is what gets printed.
	Path    string
	Pkg     string // import path of the package this file belongs to
	Syntax  *ast.File
	imports map[string]string // local name -> import path
	consts  map[string]int64  // package-level integer constants, filled per package
}

// isTest reports whether this file is a _test.go.
//
// The checks about a seeded value - a menu sort, a config value, a menu id, a
// soft-delete shape - are all about what reaches a real database through a
// migration, and a test fixture reaches none. Worse, each of those guards
// needs a test that writes the very value it rejects, so scanning test files
// makes every such guard report its own test. The import and alias checks do
// not skip tests: those are about the dependency graph, where a test file's
// import is as real as any other.
func (f *sourceFile) isTest() bool {
	return strings.HasSuffix(f.Path, "_test.go")
}

// snapshot is every Go file under the root, parsed once and shared by all the
// checks.
type snapshot struct {
	Root       string
	ModulePath string
	Fset       *token.FileSet
	Files      []*sourceFile
}

// load parses every Go file under root.
//
// Type checking is deliberately not done. Every rule here is about a literal
// written into a seed or an import that should not be there, and both are in
// the syntax; a type checker would drag in a resolvable build of the whole
// module, which is exactly what a check meant to run on a half-broken tree
// cannot depend on.
func load(root string) (*snapshot, error) {
	modulePath, err := readModulePath(root)
	if err != nil {
		return nil, err
	}

	s := &snapshot{Root: root, ModulePath: modulePath, Fset: token.NewFileSet()}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != root && skippedDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		f, err := parser.ParseFile(s.Fset, path, nil, 0)
		if err != nil {
			// A file that does not parse is someone's work in progress, not a
			// silent failure. Reporting it here would only repeat what the
			// compiler already says louder.
			return nil
		}
		dir := filepath.ToSlash(filepath.Dir(rel))
		pkg := modulePath
		if dir != "." {
			pkg = modulePath + "/" + dir
		}
		s.Files = append(s.Files, &sourceFile{
			Path:    rel,
			Pkg:     pkg,
			Syntax:  f,
			imports: fileImports(f),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.resolveConstants()
	return s, nil
}

func readModulePath(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", fmt.Errorf("%s has no module line", filepath.Join(root, "go.mod"))
}

func fileImports(f *ast.File) map[string]string {
	out := make(map[string]string, len(f.Imports))
	for _, spec := range f.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := path[strings.LastIndex(path, "/")+1:]
		if spec.Name != nil {
			name = spec.Name.Name
		}
		out[name] = path
	}
	return out
}

// resolveConstants collects package-level integer constants per directory, so a
// seed written as MenuId: demoMenuId can be compared against one written as
// MenuId: 9000. Without this the menu-id check would see the well-written code
// and miss exactly the collisions it exists to find.
func (s *snapshot) resolveConstants() {
	byPkg := map[string]map[string]int64{}
	for _, sf := range s.Files {
		consts, ok := byPkg[sf.Pkg]
		if !ok {
			consts = map[string]int64{}
			byPkg[sf.Pkg] = consts
		}
		for _, decl := range sf.Syntax.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range vs.Names {
					if i >= len(vs.Values) {
						continue
					}
					if v, ok := intLiteral(vs.Values[i]); ok {
						consts[name.Name] = v
					}
				}
			}
		}
	}
	for _, sf := range s.Files {
		sf.consts = byPkg[sf.Pkg]
	}
}

// Pos turns a token position into the file/line/col a Finding carries.
func (s *snapshot) Pos(sf *sourceFile, p token.Pos) (string, int, int) {
	pos := s.Fset.Position(p)
	return sf.Path, pos.Line, pos.Column
}

// Imports reports whether the file imports path.
func (sf *sourceFile) Imports(path string) bool {
	for _, p := range sf.imports {
		if p == path {
			return true
		}
	}
	return false
}

// ImportSpec returns the import declaration for path, for reporting a position
// on the import line itself.
func (sf *sourceFile) ImportSpec(path string) *ast.ImportSpec {
	for _, spec := range sf.Syntax.Imports {
		if p, err := strconv.Unquote(spec.Path.Value); err == nil && p == path {
			return spec
		}
	}
	return nil
}
