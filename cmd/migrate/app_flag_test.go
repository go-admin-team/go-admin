package migrate

import (
	"strings"
	"testing"

	"go-admin/cmd/migrate/migration"
)

// A mistyped --app used to be indistinguishable from an up-to-date database on
// all three paths: `migrate` printed that the app was unknown and exited 0,
// while `--dry-run` and `status` printed "nothing to apply" and "none
// recorded" - the same words a database with nothing pending produces. An
// operator scripting `migrate --app crmm && deploy` therefore deployed against
// a database the migrations never touched.
func TestAppRegistrationErrorRejectsAnUnknownCode(t *testing.T) {
	restore := appCode
	t.Cleanup(func() { appCode = restore })

	appCode = "doesnotexist"
	err := appRegistrationError()
	if err == nil {
		t.Fatal("an unregistered app code must be an error, not an empty run")
	}
	if !strings.Contains(err.Error(), `"doesnotexist"`) {
		t.Errorf("the message must quote what was typed; got %q", err)
	}
	// Listing what is registered is what turns the error into a fix: the typo
	// is usually one letter away from something in this list.
	if !strings.Contains(err.Error(), migration.FrameworkAppCode) {
		t.Errorf("the message must list the registered codes; got %q", err)
	}
}

func TestAppRegistrationErrorAcceptsWhatIsRegistered(t *testing.T) {
	restore := appCode
	t.Cleanup(func() { appCode = restore })

	for _, code := range []string{
		"",                         // no --app at all: every migration runs
		migration.FrameworkAppCode, // "core", the framework's own
		strings.ToUpper(migration.FrameworkAppCode), // codes normalize to lower case
	} {
		appCode = code
		if err := appRegistrationError(); err != nil {
			t.Errorf("appCode %q must be accepted; got %v", code, err)
		}
	}
}
