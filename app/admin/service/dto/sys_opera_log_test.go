package dto

import (
	"testing"

	"go-admin/common/global"
)

// The values moved to common/global so common/middleware would stop importing
// this package; these two names stayed behind as aliases, misspelling and all,
// because forks import them.
//
// If they ever drift apart, rows written through the two spellings land in
// different buckets and the operation-log filter silently misses half of them.
func TestDeprecatedStatusAliasesStillMatch(t *testing.T) {
	if OperaStatusEnabel != global.OperaStatusEnabled {
		t.Errorf("OperaStatusEnabel = %q, global.OperaStatusEnabled = %q",
			OperaStatusEnabel, global.OperaStatusEnabled)
	}
	if OperaStatusDisable != global.OperaStatusDisabled {
		t.Errorf("OperaStatusDisable = %q, global.OperaStatusDisabled = %q",
			OperaStatusDisable, global.OperaStatusDisabled)
	}
}
