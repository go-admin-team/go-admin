package router

import (
	"testing"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// genWritingRoutes are the three that do not read. They write Go and Vue
// source onto the host, a migration, and rows in sys_menu.
var genWritingRoutes = []string{
	"/api/v1/gen/toproject/:tableId",
	"/api/v1/gen/apitofile/:tableId",
	"/api/v1/gen/todb/:tableId",
}

// genReadingRoutes are the rest of the generator's surface. Gating the three
// above must not cost any of these: a deployment that cannot list its tables
// or preview a template has lost the feature, not secured it.
var genReadingRoutes = []string{
	"/api/v1/gen/preview/:tableId",
	"/api/v1/gen/tabletree",
	"/api/v1/db/tables/page",
	"/api/v1/db/columns/page",
}

// The endpoints that write are registered where the mode says development and
// nowhere else.
//
// They are in CasbinExclude, so Enforce never runs for them and any account
// that can log in may call them. dev is the shipped default and is where the
// generator is meant to be used; demo keeps them because DemoEvn refuses these
// three by name and saying so is the thing the demo is for. Everything else,
// including the empty mode a process gets when nothing set one, is refused by
// not existing.
func TestGeneratorWritingRoutesExistOnlyWhereTheModeAllowsIt(t *testing.T) {
	for _, tc := range []struct {
		mode     string
		expected bool
		why      string
	}{
		{"dev", true, "the shipped default, and where the generator is used"},
		{"demo", true, "registered so demo mode can refuse them by name"},
		{"test", false, "a deployment, however much it is called a test"},
		{"prod", false, "a deployment"},
		{"", false, "no mode configured is not a reason to trust the caller"},
	} {
		t.Run("mode="+tc.mode, func(t *testing.T) {
			routes := registeredRoutes(t, tc.mode)
			for _, writing := range genWritingRoutes {
				if got := routes[writing]; got != tc.expected {
					t.Errorf("mode %q: %s registered = %v, want %v (%s)",
						tc.mode, writing, got, tc.expected, tc.why)
				}
			}
		})
	}
}

// The other direction. Refusing too much is as much of a defect as refusing
// too little, and the read-only half of the generator is what a demo shows.
func TestGeneratorReadingRoutesExistInEveryMode(t *testing.T) {
	for _, mode := range []string{"dev", "demo", "test", "prod", ""} {
		t.Run("mode="+mode, func(t *testing.T) {
			routes := registeredRoutes(t, mode)
			for _, reading := range genReadingRoutes {
				if !routes[reading] {
					t.Errorf("mode %q: %s is not registered - the gate took a route that only reads",
						mode, reading)
				}
			}
		})
	}
}

// GenWriteRoutesEnabled is what cmd/api reads to decide whether to warn at
// start-up. If it and the registration ever disagree, the log says one thing
// and the engine does another, so they are checked against each other rather
// than each against a list.
func TestGenWriteRoutesEnabledAgreesWithWhatWasRegistered(t *testing.T) {
	for _, mode := range []string{"dev", "demo", "test", "prod", ""} {
		t.Run("mode="+mode, func(t *testing.T) {
			routes := registeredRoutes(t, mode)

			// registeredRoutes puts the mode back before it returns, so ask
			// the predicate under a mode set here - about the same value the
			// engine was just built under.
			previous := config.ApplicationConfig.Mode
			t.Cleanup(func() { config.ApplicationConfig.Mode = previous })
			config.ApplicationConfig.Mode = mode

			claimed := GenWriteRoutesEnabled()
			actual := routes["/api/v1/gen/todb/:tableId"]
			if claimed != actual {
				t.Errorf("mode %q: GenWriteRoutesEnabled() = %v but the route was registered = %v",
					mode, claimed, actual)
			}
		})
	}
}

// The helper restores the mode before it returns, so nothing it was asked
// about leaks into what the caller does next.
//
// Worth a test of its own because the failure is silent: a helper that left
// the mode set would make every assertion after the call read a value the
// caller did not choose, and each of those assertions would still pass for as
// long as the leaked value happened to be the right one.
func TestRegisteredRoutesRestoresTheModeBeforeReturning(t *testing.T) {
	const sentinel = "not-a-mode"

	previous := config.ApplicationConfig.Mode
	t.Cleanup(func() { config.ApplicationConfig.Mode = previous })
	config.ApplicationConfig.Mode = sentinel

	registeredRoutes(t, "prod")

	if got := config.ApplicationConfig.Mode; got != sentinel {
		t.Errorf("mode after the helper returned = %q, want %q - it was left set to what "+
			"the helper was asked about", got, sentinel)
	}
}

// Changing the mode after the routes were built unregisters nothing.
//
// buildRouter has one call site, in run(), and route registration is not on
// any phase or reload callback - so a configuration reload moves
// config.ApplicationConfig.Mode without moving the routes. From that moment
// GenWriteRoutesEnabled answers about a mode the engine was not built under.
//
// That gap is why the start-up warning tells the reader to restart rather than
// only to change the mode. This pins it: if registration ever becomes dynamic,
// this test fails and the message it justifies has to be revisited.
func TestChangingTheModeDoesNotUnregisterWhatWasAlreadyBuilt(t *testing.T) {
	built := registeredRoutes(t, "dev")
	if !built["/api/v1/gen/todb/:tableId"] {
		t.Fatal("built under dev without the writing routes, so this test asserts nothing")
	}

	previous := config.ApplicationConfig.Mode
	t.Cleanup(func() { config.ApplicationConfig.Mode = previous })
	config.ApplicationConfig.Mode = "prod"

	if GenWriteRoutesEnabled() {
		t.Fatal("the predicate still allows prod, so the disagreement below is not the one meant")
	}
	if !built["/api/v1/gen/todb/:tableId"] {
		t.Error("the route left the engine when the mode changed - registration has become " +
			"dynamic, and the start-up warning's advice to restart is now wrong")
	}
}
