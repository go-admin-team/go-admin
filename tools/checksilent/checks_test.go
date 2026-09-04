package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture writes a miniature repository and returns its root. Every check gets
// one of these carrying the exact mistake it exists to find, so a check that
// stops working fails a test rather than going quiet - which would be the same
// failure mode the tool is about.
func fixture(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	files["go.mod"] = "module go-admin\n\ngo 1.26\n"
	for name, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func check(t *testing.T, root string, opt options) []Finding {
	t.Helper()
	s, err := load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	findings, err := runChecks(s, opt)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return findings
}

func only(t *testing.T, findings []Finding, name string) []Finding {
	t.Helper()
	var out []Finding
	for _, f := range findings {
		if f.Check == name {
			out = append(out, f)
		}
	}
	return out
}

func requireOne(t *testing.T, findings []Finding, name string) Finding {
	t.Helper()
	got := only(t, findings, name)
	if len(got) != 1 {
		t.Fatalf("%s produced %d findings, want 1:\n%v", name, len(got), findings)
	}
	return got[0]
}

const frozenModelsPkg = `package models

import (
	"time"

	"gorm.io/gorm"
)

type ModelTime struct {
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type SysMenu struct {
	MenuId    int
	MenuName  string
	Component string
	MenuType  string
	Sort      int
	ModelTime
}

func (SysMenu) TableName() string { return "sys_menu" }

type SysConfig struct {
	ConfigKey   string
	ConfigValue string
	ModelTime
}

func (SysConfig) TableName() string { return "sys_config" }
`

// ---------------------------------------------------------------------------

func TestModelTimeMixDetectsARuntimeModelOnTheFrozenShape(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"app/shop/models/product.go": `package models

import frozen "go-admin/cmd/migrate/migration/models"

type Product struct {
	Id int
	frozen.ModelTime
}

func (Product) TableName() string { return "shop_product" }
`,
	})

	f := requireOne(t, check(t, root, options{}), checkModelTimeMix)
	if f.Severity != "ERROR" {
		t.Errorf("severity = %s", f.Severity)
	}
	if !strings.Contains(f.Message, "shop_product") {
		t.Errorf("message = %s", f.Message)
	}
	if f.File != "app/shop/models/product.go" {
		t.Errorf("file = %s", f.File)
	}
}

func TestModelTimeMixDetectsAPostConversionMigrationOnTheFrozenPackage(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version-local/1786700009000_seed.go": `package version_local

import "go-admin/cmd/migrate/migration/models"

func seed() interface{} { return &models.SysMenu{} }
`,
	})

	f := requireOne(t, check(t, root, options{}), checkModelTimeMix)
	if !strings.Contains(f.Message, "1786700009000") {
		t.Errorf("message = %s", f.Message)
	}
	if f.Line != 3 {
		t.Errorf("expected the import line, got line %d", f.Line)
	}
}

// A migration ordered before the conversion is right to use that package - the
// nullable column is the shape it had at the time.
func TestModelTimeMixLeavesPreConversionMigrationsAlone(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_seed.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() interface{} { return &models.SysMenu{} }
`,
	})
	if got := only(t, check(t, root, options{}), checkModelTimeMix); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// ---------------------------------------------------------------------------

func TestMenuSortOverflowIsDetectedThroughAnElidedSliceLiteral(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_seed.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() []models.SysMenu {
	return []models.SysMenu{
		{MenuId: 9000, Sort: 100},
		{MenuId: 9001, Sort: 900},
	}
}
`,
	})

	f := requireOne(t, check(t, root, options{}), checkMenuSort)
	if f.Severity != "ERROR" || !strings.Contains(f.Message, "900") {
		t.Errorf("finding = %+v", f)
	}
	if f.Line != 8 {
		t.Errorf("line = %d, want the overflowing element", f.Line)
	}
}

// The SysMenu of the service layer shares the name and has nothing to do with
// the column. Reporting it would be the false positive that gets the tool
// switched off.
func TestMenuSortIgnoresSameNamedTypesFromOtherPackages(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"app/admin/service/sys_menu.go": `package service

type SysMenu struct{ Sort int }
`,
		"app/admin/apis/sys_menu.go": `package apis

import "go-admin/app/admin/service"

func handler() interface{} { return service.SysMenu{Sort: 9000} }
`,
	})
	if got := only(t, check(t, root, options{}), checkMenuSort); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// ---------------------------------------------------------------------------

func TestConfigValueTruncationIsDetected(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_seed.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() models.SysConfig {
	return models.SysConfig{ConfigKey: "k", ConfigValue: "` + strings.Repeat("a", 300) + `"}
}
`,
	})

	f := requireOne(t, check(t, root, options{}), checkConfigValue)
	if f.Severity != "ERROR" || !strings.Contains(f.Message, "300 characters") {
		t.Errorf("finding = %+v", f)
	}
}

// varchar(255) counts characters. Counting bytes would report a Chinese value
// that fits perfectly well.
func TestConfigValueCountsRunesNotBytes(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_seed.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() models.SysConfig {
	return models.SysConfig{ConfigValue: "` + strings.Repeat("中", 200) + `"}
}
`,
	})
	if got := only(t, check(t, root, options{}), checkConfigValue); len(got) != 0 {
		t.Errorf("200 Chinese characters fit in varchar(255) but were reported: %v", got)
	}
}

// ---------------------------------------------------------------------------

func TestMenuIDCollisionAcrossFilesIsDetectedThroughConstants(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_crm.go": `package version

import "go-admin/cmd/migrate/migration/models"

const crmMenuId = 9000

func seedCrm() models.SysMenu { return models.SysMenu{MenuId: crmMenuId, MenuName: "CRM"} }
`,
		"cmd/migrate/migration/version/1786700002000_oms.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seedOms() models.SysMenu { return models.SysMenu{MenuId: 9000, MenuName: "OMS"} }
`,
	})

	f := requireOne(t, check(t, root, options{}), checkMenuIDConflict)
	if f.Severity != "ERROR" || !strings.Contains(f.Message, "9000") {
		t.Errorf("finding = %+v", f)
	}
	// The second site has to be printed too, or the report says a collision
	// happened without saying with what.
	if len(f.Related) != 1 || !strings.Contains(f.Related[0], "1786700002000_oms.go") {
		t.Errorf("related = %v", f.Related)
	}
}

// The seeds write a menu and then refer to it again to attach permissions. That
// is one menu, not two modules fighting over an id.
func TestMenuIDRepeatedWithinOneFileIsNotACollision(t *testing.T) {
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_crm.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() []models.SysMenu {
	dir := models.SysMenu{MenuId: 9000, MenuName: "CRM"}
	again := models.SysMenu{MenuId: 9000, MenuName: "CRM"}
	return []models.SysMenu{dir, again}
}
`,
	})
	if got := only(t, check(t, root, options{}), checkMenuIDConflict); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// ---------------------------------------------------------------------------

func TestContractImportBoundaryIsDetected(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/service/dto/log.go": `package dto

const OperaStatusEnabel = "1"
`,
		"common/middleware/logger.go": `package middleware

import "go-admin/app/admin/service/dto"

func status() string { return dto.OperaStatusEnabel }
`,
	})

	f := requireOne(t, check(t, root, options{}), checkImportBoundary)
	if f.Severity != "ERROR" {
		t.Errorf("severity = %s", f.Severity)
	}
	if !strings.Contains(f.Message, "go-admin/app/admin/service/dto") {
		t.Errorf("message = %s", f.Message)
	}
	if f.File != "common/middleware/logger.go" || f.Line != 3 {
		t.Errorf("position = %s:%d", f.File, f.Line)
	}
}

// A fork that drops app/admin should be able to run the tests too, so a
// test-only import is the same violation.
func TestContractImportBoundaryCoversTestFiles(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/models/user.go": "package models\n\ntype SysUser struct{}\n",
		"common/actions/permission_test.go": `package actions

import "go-admin/app/admin/models"

var _ = models.SysUser{}
`,
	})
	if got := only(t, check(t, root, options{}), checkImportBoundary); len(got) != 1 {
		t.Errorf("findings = %v", got)
	}
}

func TestContractImportBoundaryAllowsAppToImportCommon(t *testing.T) {
	root := fixture(t, map[string]string{
		"common/models/by.go": "package models\n\ntype ControlBy struct{}\n",
		"app/demo/models/product.go": `package models

import "go-admin/common/models"

type Product struct{ models.ControlBy }
`,
	})
	if got := only(t, check(t, root, options{}), checkImportBoundary); len(got) != 0 {
		t.Errorf("the dependency direction that is allowed was reported: %v", got)
	}
}

// ---------------------------------------------------------------------------

func menuNameFixture(t *testing.T, menuName, componentName string) (string, string) {
	t.Helper()
	root := fixture(t, map[string]string{
		"cmd/migrate/migration/models/models.go": frozenModelsPkg,
		"cmd/migrate/migration/version/1786700001000_seed.go": `package version

import "go-admin/cmd/migrate/migration/models"

func seed() models.SysMenu {
	return models.SysMenu{MenuId: 9001, MenuName: "` + menuName + `", MenuType: "C", Component: "/demo/product/index"}
}
`,
	})
	ui := t.TempDir()
	path := filepath.Join(ui, "views", "demo", "product", "index.vue")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	vue := "<script setup>\ndefineOptions({ name: '" + componentName + "' })\n</script>\n<template><div/></template>\n"
	if err := os.WriteFile(path, []byte(vue), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, ui
}

func TestMenuNameMismatchIsDetectedAsAWarning(t *testing.T) {
	root, ui := menuNameFixture(t, "DemoProduct", "Product")

	f := requireOne(t, check(t, root, options{UIDir: ui}), checkMenuName)
	if f.Severity != "WARN" {
		t.Errorf("severity = %s; this check must not decide an exit code yet", f.Severity)
	}
	if !strings.Contains(f.Message, "DemoProduct") || !strings.Contains(f.Message, "Product") {
		t.Errorf("message = %s", f.Message)
	}
	// Acceptance 14c.
	if !strings.Contains(f.Message, "heuristic") || !strings.Contains(f.Message, "false positives are possible") {
		t.Errorf("the warning must say it is a heuristic:\n%s", f.Message)
	}
}

func TestMenuNameMatchIsSilent(t *testing.T) {
	root, ui := menuNameFixture(t, "DemoProduct", "DemoProduct")
	if got := only(t, check(t, root, options{UIDir: ui}), checkMenuName); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// Without the frontend repository there is nothing to compare against, and the
// check has to disappear rather than guess.
func TestMenuNameIsSkippedWithoutTheUIDirectory(t *testing.T) {
	root, _ := menuNameFixture(t, "DemoProduct", "Product")
	if got := only(t, check(t, root, options{}), checkMenuName); len(got) != 0 {
		t.Errorf("reported %v without -ui-dir", got)
	}
}

// A component the frontend does not carry is the "app not installed" case F2
// handles with a placeholder; this check has nothing to say about it.
func TestMenuNameIsSilentWhenTheComponentIsMissing(t *testing.T) {
	root, _ := menuNameFixture(t, "DemoProduct", "Product")
	empty := t.TempDir()
	if got := only(t, check(t, root, options{UIDir: empty}), checkMenuName); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

func TestComponentNameParsesBothVueStyles(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		want string
	}{
		{"script setup", "<script setup>\ndefineOptions({ name: 'Product' })\n</script>", "Product"},
		{"options api", "<script>\nexport default {\n  name: 'Product',\n  data() {}\n}\n</script>", "Product"},
		{"defineComponent", "<script>\nexport default defineComponent({\n  name: 'Product'\n})\n</script>", "Product"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := componentName(tc.src)
			if !ok || got != tc.want {
				t.Errorf("componentName = %q, %v", got, ok)
			}
		})
	}
	if _, ok := componentName("<script setup>\nconst a = 1\n</script>"); ok {
		t.Error("a component with no declared name must not be compared")
	}
}

// ---------------------------------------------------------------------------

const coreContractModels = "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"

// shimFixture writes one shim file declaring ControlBy from core's contract
// package, in whichever of the two forms the caller asks for.
func shimFixture(t *testing.T, decl string) string {
	t.Helper()
	return fixture(t, map[string]string{
		"common/models/by.go": "package models\n\nimport \"" + coreContractModels + "\"\n\n" + decl + "\n",
	})
}

func TestShimAliasDetectsADefinedType(t *testing.T) {
	root := shimFixture(t, "type ControlBy models.ControlBy")

	f := requireOne(t, check(t, root, options{}), checkShimAlias)
	if f.Severity != "ERROR" {
		t.Errorf("severity = %s", f.Severity)
	}
	if !strings.Contains(f.Message, "type ControlBy = models.ControlBy") {
		t.Errorf("the message must spell out the fix; got %s", f.Message)
	}
	if f.File != "common/models/by.go" || f.Line != 5 {
		t.Errorf("position = %s:%d", f.File, f.Line)
	}
}

// The counterproof for the check above: the same fixture with the one
// character that makes it correct must produce nothing. Without this the check
// could be reporting every type declaration it sees and the test above would
// still pass.
func TestShimAliasAcceptsAnAlias(t *testing.T) {
	root := shimFixture(t, "type ControlBy = models.ControlBy")
	if got := only(t, check(t, root, options{}), checkShimAlias); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// A parenthesised type block is how a shim package with more than one type
// tends to get written, and a walker that only looked at single-spec
// declarations would skip all but the first.
func TestShimAliasReadsAParenthesisedBlock(t *testing.T) {
	root := shimFixture(t, `type (
	Model     = models.Model
	ControlBy models.ControlBy
	ModelTime = models.ModelTime
)`)
	f := requireOne(t, check(t, root, options{}), checkShimAlias)
	if !strings.Contains(f.Message, "ControlBy") {
		t.Errorf("message = %s", f.Message)
	}
}

// A defined type over a package that is not core's contract namespace is
// somebody's ordinary code. The check exists for the surface core promises to
// keep stable, and reporting anything else would make it a style rule.
func TestShimAliasIgnoresOtherPackages(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/demo/models/product.go": `package models

import "go-admin/common/models"

type Product models.Model
`,
	})
	if got := only(t, check(t, root, options{}), checkShimAlias); len(got) != 0 {
		t.Errorf("reported %v", got)
	}
}

// The version is part of core's import path and changes on every major bump.
// Matching the whole path literally would turn the check off on that day and
// say nothing about it.
func TestShimAliasSurvivesACoreMajorVersionBump(t *testing.T) {
	root := fixture(t, map[string]string{
		"common/models/by.go": `package models

import "github.com/go-admin-team/go-admin-core/v9/sdk/contract/models"

type ControlBy models.ControlBy
`,
	})
	if got := only(t, check(t, root, options{}), checkShimAlias); len(got) != 1 {
		t.Errorf("findings = %v", got)
	}
}

// A tree with no shims in it is the state of this repository until the
// contract packages are lowered, and the check saying nothing there must not
// be reported as a boundary being guarded.
func TestShimAliasCoverageIsReportedAsZeroWhenThereAreNoShims(t *testing.T) {
	root := fixture(t, map[string]string{
		"common/models/by.go": "package models\n\ntype ControlBy struct{}\n",
	})
	s, err := load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if n := ScannedShimAliases(s); n != 0 {
		t.Errorf("ScannedShimAliases = %d, want 0", n)
	}

	var buf strings.Builder
	if _, err := run(&buf, root, options{}, false); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(buf.String(), "guarded nothing") {
		t.Errorf("the summary must say the check covered nothing; got:\n%s", buf.String())
	}
}

func TestShimAliasCoverageCountsTheAliasesItGuards(t *testing.T) {
	root := shimFixture(t, `type (
	Model     = models.Model
	ControlBy = models.ControlBy
)`)
	s, err := load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if n := ScannedShimAliases(s); n != 2 {
		t.Errorf("ScannedShimAliases = %d, want 2", n)
	}
}
