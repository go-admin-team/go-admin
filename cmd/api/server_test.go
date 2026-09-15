package api

import (
	"strings"
	"testing"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
)

// freshRuntime hands the test its own Runtime and puts the old one back.
//
// Both registries close permanently the first time they are run, and
// sdk.Runtime is a process-wide singleton, so a test that runs the startup
// hooks would otherwise leave every later test in this binary registering into
// a closed registry - which is only an ERROR log, not a failure. The symptom
// is a test that passes alone and loses its routes when run with the others.
func freshRuntime(t *testing.T) {
	t.Helper()
	previous := sdk.Runtime
	t.Cleanup(func() { sdk.Runtime = previous })
	sdk.Runtime = runtime.NewConfig()
}

// Acceptance 1 and 2 together: the package-level slice a fork appends to and
// the core registry a module registers through both run, package-level first,
// registration order preserved inside each.
//
// The order matters beyond neatness. A module that appends to AppRouters has to
// import go-admin/cmd/api, which is why every module used to need a seven-line
// file in the command package; SetAppRouters is the way out of that. Running
// the old registry first is what makes the change invisible to anyone who never
// takes it.
func TestRunStartupHooksRunsBothRegistriesInOrder(t *testing.T) {
	freshRuntime(t)

	savedPackage := AppRouters
	t.Cleanup(func() { AppRouters = savedPackage })

	var order []string
	AppRouters = []func(){
		func() { order = append(order, "package-1") },
		func() { order = append(order, "package-2") },
	}
	sdk.Runtime.SetAppRouters(func() { order = append(order, "runtime-1") })
	sdk.Runtime.SetAppRouters(func() { order = append(order, "runtime-2") })

	runStartupHooks()

	const want = "package-1,package-2,runtime-1,runtime-2"
	if got := strings.Join(order, ","); got != want {
		t.Errorf("ran %q, want %q", got, want)
	}
}

// Acceptance 17: a before callback registered through core actually runs.
//
// It did not, ever: core stored the callbacks and nothing executed them, so
// SetBefore was accepted and silently did nothing. The gap survived because
// core offered the registry without ever running it, leaving each consumer to
// write - or forget - its own loop.
func TestBeforeCallbacksRun(t *testing.T) {
	freshRuntime(t)

	savedPackage := AppRouters
	t.Cleanup(func() { AppRouters = savedPackage })
	AppRouters = nil

	var order []string
	sdk.Runtime.SetBefore(func() { order = append(order, "before-1") })
	sdk.Runtime.SetBefore(func() { order = append(order, "before-2") })
	sdk.Runtime.SetAppRouters(func() { order = append(order, "router") })

	runStartupHooks()

	// Routers first, then before: both happen ahead of ListenAndServe, and a
	// router callback is what puts the engine in place for anything that comes
	// after it.
	const want = "router,before-1,before-2"
	if got := strings.Join(order, ","); got != want {
		t.Errorf("ran %q, want %q", got, want)
	}
}

// A panicking module must not take the server down with it. The guard lives in
// core; this asserts that go-admin actually goes through it rather than around
// it with a loop of its own.
func TestAPanickingRouterDoesNotStopStartup(t *testing.T) {
	freshRuntime(t)

	savedPackage := AppRouters
	t.Cleanup(func() { AppRouters = savedPackage })
	AppRouters = nil

	var order []string
	sdk.Runtime.SetAppRouters(func() { order = append(order, "first") })
	sdk.Runtime.SetAppRouters(func() { panic("a third-party module blew up") })
	sdk.Runtime.SetAppRouters(func() { order = append(order, "third") })
	sdk.Runtime.SetBefore(func() { order = append(order, "before") })

	runStartupHooks()

	const want = "first,third,before"
	if got := strings.Join(order, ","); got != want {
		t.Errorf("ran %q, want %q", got, want)
	}
}

// The default AppRouters must keep the admin routes on it. Emptying the slice
// would not fail to compile anywhere - it would just serve a server with no
// admin API and no error.
func TestAdminRouterIsRegisteredOnThePackageSlice(t *testing.T) {
	if len(AppRouters) == 0 {
		t.Fatal("AppRouters is empty; the admin router is registered in init()")
	}
}
