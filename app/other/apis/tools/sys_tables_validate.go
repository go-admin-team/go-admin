package tools

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"go-admin/app/other/models/tools"
)

// jsonFieldPattern mirrors genInfoForm.vue's businessName rule
// (`/^[a-z][A-Za-z]+$/`) - jsonField has never had a format rule of its own,
// unlike businessName/tableName/className, and API契约.md §1.2 recommends
// tightening it to the same identifier shape the other three already use.
var jsonFieldPattern = regexp.MustCompile(`^[a-z][A-Za-z]+$`)

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
