package tools

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"go-admin/app/other/models/tools"
)

// jsonFieldPattern accepts any legal JS/TS identifier that starts with a
// lowercase letter - not businessName's rule.
//
// This used to be businessName's own pattern (^[a-z][A-Za-z]+$, requiring at
// least two letters and no digits), copied over on the theory that jsonField
// "should tighten to the same identifier shape". That theory does not hold:
// businessName is typed by a person on genInfoForm.vue, so a strict pattern
// is a reasonable guardrail on human input. jsonField is computed by the
// importer from the column name (sys_tables.go's namelist/JsonField loop) -
// nobody types it, so the same pattern only rejects names the importer
// legitimately produces. A one-letter column ("x") or a column ending in a
// digit ("address2", "a1") both import to a single camelCase word with no
// separators to re-capitalize, and both used to fail this check - meaning a
// table that merely contained such a column could never save any config
// again, unrelated columns included, since this check runs over every
// column on every Update.
//
// What still has to be rejected is a jsonField that cannot be a raw object
// key at all: empty, containing whitespace/punctuation, or leading with a
// digit (`2faEnabled: 1` is not valid JS - identifiers cannot start with a
// digit, and this is what lands as the property name in gen.go's generated
// interface / lang file, both unquoted). Hence still anchoring on a
// lowercase letter first, but no longer requiring a second character or
// forbidding digits after it.
var jsonFieldPattern = regexp.MustCompile(`^[a-z][A-Za-z0-9]*$`)

// colWidthMin/colWidthMax are API契约.md §2.1's suggested range for colWidth.
const (
	colWidthMin = 40
	colWidthMax = 800
)

// expressionMarkers flags the "meant to be evaluated" shapes API契约.md §2.1
// says defaultValue must not carry: it is spliced into the generated
// defaultModel() as a literal and never evaluated, so anything that looks
// like a function call or a block is rejected outright rather than
// generating code that silently does nothing.
var expressionMarkers = []string{"(", ")", "{", "}", "`", ";", "=>"}

// validateAndSanitizeColumns enforces PRD 010 F10 on the columns carried by
// a table update (sys_tables.go:357's Update handler, the one bind-and-save
// path with no field-level validation at all - see API契约.md §1.2/§2.1,
// decision D6).
//
// jsonField and defaultValue problems reject the request outright: letting
// either through would corrupt the generated i18n file silently (a
// duplicate or malformed jsonField becomes a duplicate or invalid key in
// gen/{PackageName}/{BusinessName}.ts, see the lang-zh/lang-en templates).
// An out-of-range colWidth does not reject - §2.1 says it "falls back to
// the inferred value", so this resets it to the 0 sentinel in place and lets
// R2's inference take over, the same as if the field had never been set.
func validateAndSanitizeColumns(columns []tools.SysColumns) error {
	seen := make(map[string]bool, len(columns))
	for i := range columns {
		col := &columns[i]

		if !jsonFieldPattern.MatchString(col.JsonField) {
			return fmt.Errorf("jsonField 格式不合法：%q，须以小写字母开头且只能包含英文字母", col.JsonField)
		}
		if seen[col.JsonField] {
			return fmt.Errorf("jsonField 在同一张表内重复：%q", col.JsonField)
		}
		seen[col.JsonField] = true

		if col.ColWidth != 0 && (col.ColWidth < colWidthMin || col.ColWidth > colWidthMax) {
			col.ColWidth = 0
		}

		for _, marker := range expressionMarkers {
			if strings.Contains(col.DefaultValue, marker) {
				return fmt.Errorf("defaultValue 不允许包含表达式或函数调用内容：%q", col.DefaultValue)
			}
		}
	}
	return nil
}

// validateBusinessNameUnique enforces PRD 010 F10's other half: two tables
// sharing (packageName, businessName) write the same generated language
// pack path, gen/{PackageName}/{BusinessName}.ts (see gen.go's
// NOActionsGen), so the second one silently overwrites the first's
// translations. tableID excludes the row being saved, so a table updating
// its own unchanged name does not trip the check on itself.
//
// G10's other concern - colliding with the built-in admin/* i18n namespace -
// does not apply here anymore: D9 moved generated keys to their own gen/
// namespace, so this only has to guard generated tables against each other.
func validateBusinessNameUnique(db *gorm.DB, packageName, businessName string, tableID int) error {
	var count int64
	err := db.Table("sys_tables").
		Where("package_name = ? AND business_name = ? AND table_id != ?", packageName, businessName, tableID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("packageName=%q 下 businessName=%q 已被其它表使用", packageName, businessName)
	}
	return nil
}

// The three fields gen.go joins into the paths it writes to. Each is checked
// against what it has to be where it lands, not against one shared pattern:
//
//   - packageName names a Go package and the app/{packageName} directory.
//     Lowercase letters and digits, as a Go package name should be; no
//     hyphen, which Go rejects, and no underscore, which genInfoForm.vue
//     already refuses.
//   - tableName is the imported table's own name, and becomes a .go file name
//     and, with "_" turned into "-", a .ts/.vue one. Letters, digits and
//     underscores: what a table the importer can read is normally called, and
//     nothing that can step out of a directory.
//   - businessName is a JavaScript identifier and a .ts file name. It starts
//     with a lowercase letter, as genInfoForm.vue requires, but may carry
//     digits, which a name derived from a table such as order2 does.
//
// None of the three allows a dot or a path separator, so none can name a
// parent directory or an absolute path. The rules are no stricter than the
// config page's own, so nothing it accepts is refused here.
var (
	packageNamePattern  = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	tableNamePattern    = regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	businessNamePattern = regexp.MustCompile(`^[a-z][A-Za-z0-9]*$`)
)

// validateGenPathFields checks the fields gen.go builds file paths from. It
// runs where a configuration is saved and again before files are written, so
// a row saved before this check existed is refused at generation rather than
// trusted because it is already in the database.
func validateGenPathFields(tab tools.SysTables) error {
	if !packageNamePattern.MatchString(tab.PackageName) {
		return fmt.Errorf("packageName=%q 不合法：只能包含小写字母和数字，且以字母开头", tab.PackageName)
	}
	if !tableNamePattern.MatchString(tab.TBName) {
		return fmt.Errorf("tableName=%q 不合法：只能包含字母、数字和下划线", tab.TBName)
	}
	if !businessNamePattern.MatchString(tab.BusinessName) {
		return fmt.Errorf("businessName=%q 不合法：只能包含字母和数字，且以小写字母开头", tab.BusinessName)
	}
	return nil
}
