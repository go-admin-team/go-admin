package api

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
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
