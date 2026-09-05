package api

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
)

// freePort returns a port nothing is listening on. It is inherently a guess -
// the port is free when it is handed back and could be taken a moment later -
// but every alternative needs the caller to hold the listener, which is the one
// thing these tests cannot do.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("probe listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return port
}

// AfterListen promises a hook that the port is reachable. Both halves of that
// are asserted here, and in one test rather than two, because the phase seals
// itself once it has run: a second test calling RunPhase again would find a
// closed registry and pass while proving nothing.
//
// The failing bind comes first for the same reason. It must leave the phase
// unsealed, which is only visible if nothing has sealed it yet.
func TestAfterListenIsAnnouncedOnlyOnceThePortIsBound(t *testing.T) {
	// The pause makes the "announced synchronously" claim testable: if the
	// announcement were moved onto a goroutine, startServing would return
	// while the hook was still sleeping and the count below would be zero.
	var ran int
	sdk.Runtime.SetPhase(runtime.AfterListen, func() {
		time.Sleep(50 * time.Millisecond)
		ran++
	})

	// Somebody else already has the port. Under ListenAndServe this surfaced
	// on the serving goroutine, far too late to stop the announcement.
	taken, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("occupy: %v", err)
	}
	defer func() { _ = taken.Close() }()

	blocked := &http.Server{Addr: taken.Addr().String(), Handler: http.NewServeMux()}
	if err := startServing(blocked, false, "", ""); err == nil {
		t.Fatal("startServing returned no error for a port that was already taken")
	}
	if ran != 0 {
		t.Errorf("AfterListen ran %d times after a failed bind; a hook there is told the port is reachable", ran)
	}
	if sdk.Runtime.PhaseSealed(runtime.AfterListen) {
		t.Error("a failed bind sealed AfterListen, so the phase could never run for a server that did start")
	}

	// A certificate that cannot be read is the other way to fail before there
	// is anything to announce. ServeTLS reads it on the serving goroutine, so
	// without the check in startServing this would be a hook told the port was
	// reachable while the server was already on its way down.
	if err := startServing(&http.Server{Addr: "127.0.0.1:0"}, true, "no-such.pem", "no-such.key"); err == nil {
		t.Fatal("startServing returned no error for a certificate that does not exist")
	}
	if ran != 0 {
		t.Errorf("AfterListen ran %d times after a certificate failure", ran)
	}
	if sdk.Runtime.PhaseSealed(runtime.AfterListen) {
		t.Error("a certificate failure sealed AfterListen")
	}

	// And now a bind that works.
	port := freePort(t)
	srv := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: http.NewServeMux()}
	if err := startServing(srv, false, "", ""); err != nil {
		t.Fatalf("startServing on a free port: %v", err)
	}
	defer func() { _ = srv.Close() }()

	// Checked the instant startServing returns, so this is also the assertion
	// that it did not return early: an asynchronous announcement would still
	// be inside the sleep. Synchrony matters because an announcement that
	// overlaps the wait below could, on a fast SIGTERM, have the shutdown
	// callbacks finish before the startup ones.
	if ran != 1 {
		t.Fatalf("AfterListen ran %d times, want 1", ran)
	}

	// The claim is not "Serve was called" but "the port answers". Dial it.
	c, err := net.DialTimeout("tcp", srv.Addr, 5*time.Second)
	if err != nil {
		t.Fatalf("AfterListen ran but the port does not answer: %v", err)
	}
	_ = c.Close()
}

// BeforeRouter is the last point at which a module can still affect how routes
// are built, so it has to run while there is no engine yet. The before registry
// is a different moment despite the name: those callbacks run after initRouter
// has built the engine.
//
// The two are two lines apart in buildRouter, and calling them equivalent is a
// mistake this repository has already made in writing. Until this test the
// ordering was checked by reading - which is how the stop signals came to be
// armed after the readiness banner in the same file.
func TestBeforeRouterRunsWhileThereIsNoEngine(t *testing.T) {
	freshRuntime(t)

	// AuthInit reads these two package-level values and nothing else. No
	// database is involved in building a router: the handlers are registered,
	// not called.
	config.ApplicationConfig.Mode = "dev"
	config.JwtConfig.Secret = "test-secret-for-the-router-build"

	type observation struct {
		ran        int
		engineWas  interface{}
		engineSeen bool
	}
	var phase, before observation

	sdk.Runtime.SetPhase(runtime.BeforeRouter, func() {
		phase.ran++
		phase.engineWas = sdk.Runtime.GetEngine()
		phase.engineSeen = true
	})
	sdk.Runtime.SetBefore(func() {
		before.ran++
		before.engineWas = sdk.Runtime.GetEngine()
		before.engineSeen = true
	})

	buildRouter()

	if phase.ran != 1 {
		t.Fatalf("BeforeRouter ran %d times, want 1", phase.ran)
	}
	if !phase.engineSeen || phase.engineWas != nil {
		t.Errorf("BeforeRouter saw engine %v, want nil: it is meant to run before initRouter builds one", phase.engineWas)
	}

	if before.ran != 1 {
		t.Fatalf("the before registry ran %d times, want 1", before.ran)
	}
	if before.engineWas == nil {
		t.Error("a before callback saw no engine; that registry is meant to run after initRouter, and describing it as equivalent to BeforeRouter is the error this asserts against")
	}

	if sdk.Runtime.GetEngine() == nil {
		t.Error("buildRouter returned with no engine built")
	}
}
