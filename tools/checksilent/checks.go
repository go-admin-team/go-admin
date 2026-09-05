package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Check names. They appear in every message and in the CI log, so they are
// what people will search for.
const (
	checkModelTimeMix   = "modeltime-mix"
	checkMenuSort       = "menu-sort-overflow"
	checkMenuName       = "menu-name-mismatch"
	checkConfigValue    = "config-value-truncation"
	checkMenuIDConflict = "menu-id-collision"
	checkImportBoundary = "contract-import-boundary"
	checkShimAlias      = "contract-shim-alias"
)

// Package paths, relative to the module. Spelled once so a module rename
// touches one place.
const (
	pkgFrozenModels = "cmd/migrate/migration/models"
	pkgRuntimeModel = "common/models"
	pkgAdminModels  = "app/admin/models"
)

// options are the run-time knobs. Only the frontend directory is one: every
// other check either applies or does not, with nothing to configure.
type options struct {
	// UIDir is the go-admin-ui src directory. Empty disables checkMenuName,
	// which is the only check that needs a second repository.
	UIDir string
}

// runChecks runs every check over one parse of the tree.
func runChecks(s *snapshot, opt options) ([]Finding, error) {
	var out []Finding
	out = append(out, checkModelTimeMixing(s)...)
	out = append(out, checkMenuSortOverflow(s)...)
	out = append(out, checkConfigValueLength(s)...)
	out = append(out, checkMenuIDCollisions(s)...)
	out = append(out, checkContractImportBoundary(s)...)
	out = append(out, checkContractShimAlias(s)...)

	if opt.UIDir != "" {
		fs, err := checkMenuNames(s, opt.UIDir)
		if err != nil {
			return nil, err
		}
		out = append(out, fs...)
	}

	sortFindings(out)
	return out, nil
}

func (s *snapshot) pkg(rel string) string { return s.ModulePath + "/" + rel }

func (s *snapshot) finding(sev Severity, check string, sf *sourceFile, pos posLike, format string, args ...interface{}) Finding {
	file, line, col := s.Pos(sf, pos.Pos())
	return Finding{
		Check:    check,
		Severity: sev.String(),
		File:     file,
		Line:     line,
		Col:      col,
		Message:  fmt.Sprintf(format, args...),
		severity: sev,
	}
}

type posLike interface{ Pos() token.Pos }

// ---------------------------------------------------------------------------
// check 1: the two ModelTime flavours
// ---------------------------------------------------------------------------

// checkModelTimeMixing reports code that mixes the repository's two soft-delete
// shapes.
//
// cmd/migrate/migration/models.ModelTime declares a nullable gorm.DeletedAt;
// common/models.ModelTime declares the NOT NULL millisecond marker. Mixing them
// on one table is not a compile error and not a run-time error either: gorm
// scopes the nullable flavour as "WHERE deleted_at IS NULL" while live rows hold
// 0, so every row of the table becomes invisible and the feature reading it
// simply returns nothing. sys_columns and sys_tables sat in that state until
// 1786700004000 - the code generator listed no tables at all and reported no
// error.
//
// Two shapes are reported, and both are unambiguous:
//
//  1. a runtime model under app/ that embeds the frozen package's time struct -
//     always wrong, that package is the shape the columns had before the
//     conversion;
//  2. a migration ordered after the conversion that imports the frozen package
//     - AGENTS.md states this rule, and the version/ directory already has a
//     test for it; this extends it to version-local/, where third-party and
//     downstream migrations live and where no test was watching.
//
// Not reported: that two model packages describe the same table with different
// flavours. That is true of a dozen tables on purpose - the frozen package is
// correct for the migrations that predate the conversion - so reporting it
// would be reporting the design.
func checkModelTimeMixing(s *snapshot) []Finding {
	var out []Finding
	frozen := s.pkg(pkgFrozenModels)

	for _, sf := range s.Files {
		if sf.isTest() {
			continue
		}
		if strings.HasPrefix(sf.Path, "app/") && sf.Imports(frozen) {
			tables := tableNames(sf)
			for name, st := range structTypes(sf) {
				table, isModel := tables[name]
				if !isModel {
					continue
				}
				for _, emb := range embeddedTypes(sf, st) {
					if emb[0] != frozen {
						continue
					}
					out = append(out, s.finding(Error, checkModelTimeMix, sf, st,
						"runtime model %s (table %s) embeds %s.%s, whose DeletedAt is the nullable pre-conversion shape;\n"+
							"    gorm will query this table with deleted_at IS NULL while live rows hold 0, and it will return nothing.\n"+
							"    Embed %s.ModelTime instead.",
						name, table, pkgFrozenModels, emb[1], pkgRuntimeModel))
				}
			}
		}

		version, isMigration := migrationVersion(sf.Path)
		if !isMigration || version <= softDeleteConversion || !sf.Imports(frozen) {
			continue
		}
		spec := sf.ImportSpec(frozen)
		out = append(out, s.finding(Error, checkModelTimeMix, sf, spec,
			"migration %d is ordered after the soft-delete conversion (%d) but seeds through %s;\n"+
				"    that package writes a nullable deleted_at into a NOT NULL column, and reads through it match no rows.\n"+
				"    Use the runtime models under app/ instead.",
			version, softDeleteConversion, pkgFrozenModels))
	}
	return out
}

// softDeleteConversion is the version at which deleted_at stopped being a
// nullable timestamp and became the NOT NULL millisecond marker.
const softDeleteConversion = 1786700003000

// ---------------------------------------------------------------------------
// check 2: menu sort overflows a tinyint
// ---------------------------------------------------------------------------

// checkMenuSortOverflow reports a seeded menu sort outside a tinyint.
//
// sys_menu.sort is `gorm:"size:4"`, which MySQL builds as a tinyint holding
// -128..127. sqlite ignores the width, so an overflowing value passes every
// local test and fails on a real install - with Error 1264, partway through a
// migration that is not transactional, leaving every later migration unapplied.
// That is how a seeded Sort: 900 once stopped the run before the soft-delete
// conversion and left nobody able to log in.
func checkMenuSortOverflow(s *snapshot) []Finding {
	const (
		min = -128
		max = 127
	)
	var out []Finding
	for _, sf := range s.Files {
		if sf.isTest() {
			continue
		}
		forEachStructLiteral(sf, func(lit structLiteral) {
			if !s.isMenuModel(lit) {
				return
			}
			expr, ok := field(lit.Lit, "Sort")
			if !ok {
				return
			}
			v, ok := intValue(sf, expr)
			if !ok || (v >= min && v <= max) {
				return
			}
			out = append(out, s.finding(Error, checkMenuSort, sf, expr,
				"menu sort %d does not fit a tinyint (%d..%d);\n"+
					"    MySQL rejects it with Error 1264 and the migration stops there, leaving later migrations unapplied.",
				v, min, max))
		})
	}
	return out
}

// ---------------------------------------------------------------------------
// check 4: sys_config value truncation
// ---------------------------------------------------------------------------

// checkConfigValueLength reports a seeded sys_config value longer than the
// column.
//
// config_value is varchar(255). MySQL outside strict mode truncates rather than
// refusing, so the migration succeeds, the row is written, and the setting is
// silently half of what was intended.
//
// Counted in runes, not bytes, because varchar(255) counts characters - byte
// counting would flag Chinese values that fit.
func checkConfigValueLength(s *snapshot) []Finding {
	const limit = 255
	var out []Finding
	for _, sf := range s.Files {
		if sf.isTest() {
			continue
		}
		forEachStructLiteral(sf, func(lit structLiteral) {
			if lit.Name != "SysConfig" || !s.isModelPackage(lit.PkgPath) {
				return
			}
			expr, ok := field(lit.Lit, "ConfigValue")
			if !ok {
				return
			}
			v, ok := stringValue(sf, expr)
			if !ok {
				return
			}
			if n := len([]rune(v)); n > limit {
				out = append(out, s.finding(Error, checkConfigValue, sf, expr,
					"sys_config.config_value is %d characters, over the varchar(%d) column;\n"+
						"    MySQL outside strict mode truncates instead of failing, so the migration succeeds with half the value.",
					n, limit))
			}
		})
	}
	return out
}

// ---------------------------------------------------------------------------
// check 5: hard-coded menu ids colliding
// ---------------------------------------------------------------------------

// checkMenuIDCollisions reports the same menu id seeded from more than one file.
//
// menu_id is the primary key and every seed is an upsert, so two modules that
// pick the same id do not collide loudly - the second overwrites the first, and
// which one wins depends on migration order. One module's menu quietly becomes
// the other's.
//
// Only across files. Inside one file the same id appearing twice is the same
// menu being written and then referenced, which is how the seeds are written
// today and not a mistake.
func checkMenuIDCollisions(s *snapshot) []Finding {
	type site struct {
		sf   *sourceFile
		expr posLike
		file string
		line int
	}
	sites := map[int64][]site{}

	for _, sf := range s.Files {
		if sf.isTest() {
			continue
		}
		forEachStructLiteral(sf, func(lit structLiteral) {
			if !s.isMenuModel(lit) {
				return
			}
			expr, ok := field(lit.Lit, "MenuId")
			if !ok {
				return
			}
			v, ok := intValue(sf, expr)
			if !ok {
				return
			}
			file, line, _ := s.Pos(sf, expr.Pos())
			sites[v] = append(sites[v], site{sf: sf, expr: expr, file: file, line: line})
		})
	}

	ids := make([]int64, 0, len(sites))
	for id := range sites {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	var out []Finding
	for _, id := range ids {
		group := sites[id]
		files := map[string]bool{}
		for _, st := range group {
			files[st.file] = true
		}
		if len(files) < 2 {
			continue
		}
		sort.Slice(group, func(i, j int) bool {
			if group[i].file != group[j].file {
				return group[i].file < group[j].file
			}
			return group[i].line < group[j].line
		})
		related := make([]string, 0, len(group)-1)
		for _, st := range group[1:] {
			related = append(related, fmt.Sprintf("also at %s:%d", st.file, st.line))
		}
		f := s.finding(Error, checkMenuIDConflict, group[0].sf, group[0].expr,
			"menu id %d is seeded from %d files;\n"+
				"    menu_id is the primary key and the seeds upsert, so whichever migration runs last overwrites the other's menu.",
			id, len(files))
		f.Related = related
		out = append(out, f)
	}
	return out
}

// ---------------------------------------------------------------------------
// check 6: the contract packages must not import app/
// ---------------------------------------------------------------------------

// contractRoots are the trees that may not depend on a business module. core/
// does not exist in this repository yet and is listed because the boundary is
// declared for both in docs/contract.md; naming it here means the check is
// already in place the day the directory appears.
//
// Which of them actually exist is reported by ScannedContractRoots, because a
// root that is absent contributes nothing and a check that silently covers less
// than it claims is worse than no check: it leaves people believing a boundary
// is guarded when nothing is guarding it.
var contractRoots = []string{"common/", "core/"}

// ScannedContractRoots splits contractRoots by whether the snapshot actually
// holds files under them, so the summary can name what was covered.
func ScannedContractRoots(s *snapshot) (scanned, absent []string) {
	for _, root := range contractRoots {
		found := false
		for _, sf := range s.Files {
			if strings.HasPrefix(sf.Path, root) {
				found = true
				break
			}
		}
		if found {
			scanned = append(scanned, strings.TrimSuffix(root, "/"))
		} else {
			absent = append(absent, strings.TrimSuffix(root, "/"))
		}
	}
	return scanned, absent
}

// checkContractImportBoundary reports a contract package importing app/.
//
// docs/contract.md promises four packages under common/ as the surface an app
// may build on. A promise like that stops being true the moment the surface
// imports one particular app: a fork that replaces app/admin then cannot
// compile common/middleware, and an app can no longer be built against the
// contract alone. Nothing about it fails visibly - it fails when somebody tries
// to take the framework apart, which is the whole point of the exercise.
//
// Test files count. A fork that drops app/admin should be able to run go test
// ./... too.
func checkContractImportBoundary(s *snapshot) []Finding {
	appPrefix := s.ModulePath + "/app/"
	var out []Finding
	for _, sf := range s.Files {
		inContract := false
		for _, root := range contractRoots {
			if strings.HasPrefix(sf.Path, root) {
				inContract = true
				break
			}
		}
		if !inContract {
			continue
		}
		for _, spec := range sf.Syntax.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil || !strings.HasPrefix(path, appPrefix) {
				continue
			}
			out = append(out, s.finding(Error, checkImportBoundary, sf, spec,
				"%s is a contract package and imports %s;\n"+
					"    a fork that replaces or drops that app can then no longer compile the contract surface it was told to build on.\n"+
					"    Move what is shared down into common/, or out of the contract package.",
				sf.Path, path))
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// check 7: a shim of a core contract type must be an alias
// ---------------------------------------------------------------------------

// coreModulePrefix and coreContractSegment together identify a package under
// core's contract namespace. Matched as prefix plus segment rather than as one
// literal path so that a major-version bump of core - which rewrites the
// /v2 in every import - does not quietly turn this check off.
const (
	coreModulePrefix    = "github.com/go-admin-team/go-admin-core/"
	coreContractSegment = "/sdk/contract/"
)

// isCoreContractPkg reports whether an import path names one of core's
// contract packages.
func isCoreContractPkg(path string) bool {
	return strings.HasPrefix(path, coreModulePrefix) && strings.Contains(path, coreContractSegment)
}

// checkContractShimAlias reports a shim of a core contract type that was
// written as a defined type instead of an alias.
//
//	type ControlBy = models.ControlBy   // alias: same type, same method set
//	type ControlBy models.ControlBy     // defined type: methods are gone
//
// The two lines differ by one character and by everything else. A defined type
// takes the underlying struct and none of the methods declared on it, so a
// model embedding the second one no longer has SetCreateBy or SetUpdateBy and
// no longer satisfies ActiveRecord - which is not a warning, it is a compile
// error, but only in code that actually uses the method set.
//
// That is why the compiler is not enough on its own. This repository exercises
// some of the contract types through interfaces and some not at all; the ones
// it does not exercise compile perfectly well as defined types here and break
// in a third-party application, or in a fork's own module, which is where
// nobody is looking. The check costs one field of the AST - a type alias
// records the position of its '=' - and covers the surface uniformly rather
// than covering whatever app/demo happens to touch this month.
//
// The trigger is the right-hand side, not a list of names: any type declared
// from a core contract package is one of these, whoever wrote it and whenever
// it was added. A type declared from a local struct literal is not caught by
// this - see ScannedShimAliases, which is what stops a run over a tree with no
// shims in it from reading as a clean bill of health.
func checkContractShimAlias(s *snapshot) []Finding {
	var out []Finding
	for _, sf := range s.Files {
		forEachTypeSpec(sf, func(ts *ast.TypeSpec) {
			pkg, name, ok := qualifiedType(sf, ts.Type)
			if !ok || !isCoreContractPkg(pkg) {
				return
			}
			if ts.Assign.IsValid() {
				return // "type X = pkg.Y", which is what it must be
			}
			out = append(out, s.finding(Error, checkShimAlias, sf, ts,
				"%s is declared from %s.%s as a defined type, not an alias;\n"+
					"    a defined type keeps the fields and drops the method set, so anything embedding it stops satisfying\n"+
					"    the interfaces it satisfied before - here it may still compile, in a fork or a third-party app it does not.\n"+
					"    Write it as: type %s = %s.%s",
				ts.Name.Name, path.Base(pkg), name, ts.Name.Name, path.Base(pkg), name))
		})
	}
	return out
}

// ScannedShimAliases counts the type aliases into core's contract packages the
// snapshot holds, so the summary can say whether checkContractShimAlias found
// anything to guard at all.
//
// Reported for the same reason ScannedContractRoots is: before the contract
// packages are lowered into core there are no shims here, the check has
// nothing to look at, and a run that printed nothing would look exactly like a
// run over a tree that passed.
func ScannedShimAliases(s *snapshot) int {
	n := 0
	for _, sf := range s.Files {
		forEachTypeSpec(sf, func(ts *ast.TypeSpec) {
			pkg, _, ok := qualifiedType(sf, ts.Type)
			if ok && isCoreContractPkg(pkg) && ts.Assign.IsValid() {
				n++
			}
		})
	}
	return n
}

// forEachTypeSpec visits every type declaration in the file, including the
// ones inside a parenthesised type block.
func forEachTypeSpec(sf *sourceFile, fn func(*ast.TypeSpec)) {
	for _, decl := range sf.Syntax.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				fn(ts)
			}
		}
	}
}

// qualifiedType resolves a type expression that names a type in another
// package, returning that package's import path and the type name. A bare
// identifier, a struct literal or anything else reports false: this asks
// specifically "is the right-hand side pkg.Name", which is the shape both a
// correct shim and the mistake it guards against are written in.
func qualifiedType(sf *sourceFile, typ ast.Expr) (string, string, bool) {
	sel, ok := typ.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", "", false
	}
	path, ok := sf.imports[ident.Name]
	if !ok {
		return "", "", false
	}
	return path, sel.Sel.Name, true
}

// migrationVersion reads the 13-digit timestamp a migration file name starts
// with. Files outside the two migration directories are not migrations, however
// they are named.
func migrationVersion(rel string) (int64, bool) {
	dir := path.Dir(rel)
	if dir != "cmd/migrate/migration/version" && dir != "cmd/migrate/migration/version-local" {
		return 0, false
	}
	name := path.Base(rel)
	if len(name) < 13 {
		return 0, false
	}
	v, err := strconv.ParseInt(name[:13], 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// isMenuModel reports whether a literal describes a menu row, whichever of
// the two shapes it is written in.
//
// A host module seeds a menu by building the SysMenu model directly. An
// application installed from outside this repository cannot reach that type,
// so it describes the same row as a seed.MenuSpec and hands it to the host's
// Seeder. Both end up in sys_menu and both are subject to its column widths,
// so a check that knew only the first shape would go quiet exactly when the
// author is furthest from the schema it protects.
//
// That is not hypothetical: this repository's own reference application was
// written with a Sort of 200 - past the tinyint sys_menu.sort is built as -
// and this check passed it, because a MenuSpec is not a SysMenu.
func (s *snapshot) isMenuModel(lit structLiteral) bool {
	if lit.Name == "MenuSpec" && isCoreContractPkg(lit.PkgPath) {
		return true
	}
	return lit.Name == "SysMenu" && s.isModelPackage(lit.PkgPath)
}

func (s *snapshot) isModelPackage(path string) bool {
	return path == s.pkg(pkgFrozenModels) || path == s.pkg(pkgAdminModels)
}
