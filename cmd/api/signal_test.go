package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk"

	otherrouter "go-admin/app/other/router"
)

// The signal path cannot be exercised in-process: delivering a signal to the
// test binary would race with the test framework, and the disposition changes
// are global. So the test re-executes itself as a child, and the child runs
// gracefulShutdown - the same function run() runs, not a second copy of the
// sequence. A test that reproduces the sequence asserts against its own copy:
// move BeginDraining after the drain window and the process regresses while
// the test stays green, which is the failure mode this file exists to avoid.
//
// The child serves the real probe routes on an http.Server of its own rather
// than the configured one: this repository's CI has no database
// (.github/workflows/go.yml runs neither MySQL nor a sqlite-tagged build), and
// none of what is under test needs one. /ready answers 503 either way - with
// no database its checks fail - so the assertions below are on the draining
// answer specifically, not on the status code alone.
const (
	childEnv         = "GO_ADMIN_SIGNAL_CHILD"
	childStuckEnv    = "GO_ADMIN_SIGNAL_CHILD_STUCK"
	childHangConn    = "GO_ADMIN_SIGNAL_CHILD_HANGCONN"
	childSlowCleanup = "GO_ADMIN_SIGNAL_CHILD_SLOWCLEANUP"
	childDrainMS     = "GO_ADMIN_SIGNAL_CHILD_DRAIN_MS"
	markerAddr       = "CHILD-ADDR"
	markerReady      = "CHILD-READY"
	markerSignal     = "CHILD-SIGNAL"
	markerShutdown   = "CHILD-SHUTDOWN-OK"
	markerCleanup    = "CHILD-CLEANUP-RAN"
	markerTook       = "CHILD-TOOK-NS"
	markerExiting    = "CHILD-EXITING"
)

// childPingRoute is an ordinary route, registered beside the probes so the
// window can be checked for what it promises: requests arriving inside it are
// served, not refused. Refusing them would move the outage earlier instead of
// avoiding it.
const childPingRoute = "/signal-test-ping"

var (
	readyPath  = otherrouter.APIPrefix + otherrouter.ReadyPath
	healthPath = otherrouter.APIPrefix + otherrouter.HealthPath
	pingPath   = otherrouter.APIPrefix + childPingRoute
)

// TestSignalChild is the child process. It is skipped in a normal run.
func TestSignalChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("child process entry point")
	}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group(otherrouter.APIPrefix)
	otherrouter.RegisterMonitorRouter(v1)
	v1.GET(childPingRoute, func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		os.Exit(3)
	}
	// accepted fires once the server has taken a connection off the listener.
	// Dialling is not enough: Shutdown only waits for connections the server
	// has already accepted, so calling it between the dial and the accept
	// finds nothing to wait for and returns immediately.
	accepted := make(chan struct{}, 1)
	srv := &http.Server{
		Handler: engine,
		ConnState: func(_ net.Conn, state http.ConnState) {
			if state == http.StateNew {
				select {
				case accepted <- struct{}{}:
				default:
				}
			}
		},
	}
	go func() { _ = srv.Serve(ln) }()

	// The budget the child spends. Nothing here calls bootstrap.SetupConfig, so
	// with no environment set this is the budget of a deployment that
	// configures no extend.shutdown section at all.
	b := defaultBudget()
	if ms := os.Getenv(childDrainMS); ms != "" {
		n, err := strconv.Atoi(ms)
		if err != nil {
			fmt.Println("drain:", err)
			os.Exit(4)
		}
		b.drain = time.Duration(n) * time.Millisecond
	}

	// A BeforeExit callback, registered the way a module would. What the tests
	// below care about is whether it runs at all - after a Shutdown that
	// failed, and after its own budget has been spent.
	sdk.Runtime.SetShutdown(func(ctx context.Context) {
		switch {
		case os.Getenv(childStuckEnv) == "1":
			// Stands in for a cleanup hook that never finishes. The point of
			// restoring the signal disposition after the drain window is that
			// a second signal still reaches the default handler and kills this.
			time.Sleep(2 * time.Minute)
		case os.Getenv(childSlowCleanup) == "1":
			// Outlasts the budget on purpose, and does not consult ctx - which
			// is the case the contract is explicit about: what the context
			// bounds is the wait, not the work.
			time.Sleep(2 * time.Second)
		}
		fmt.Println(markerCleanup)
		_ = os.Stdout.Sync()
	})
	switch {
	case os.Getenv(childStuckEnv) == "1":
		b.cleanup = 2 * time.Minute
	case os.Getenv(childSlowCleanup) == "1":
		b.cleanup = 300 * time.Millisecond
	}

	// Arm before announcing readiness. Doing it the other way round leaves a
	// window in which the parent's signal reaches the default handler and
	// kills the child before any of this runs - which is exactly the failure
	// this whole change is about, so the test must not reproduce it by
	// accident.
	quit, disarm := armStopSignals()

	fmt.Println(markerAddr, ln.Addr().String())
	fmt.Println(markerReady)
	_ = os.Stdout.Sync()

	sig := <-quit
	fmt.Println(markerSignal, sig)
	_ = os.Stdout.Sync()

	if os.Getenv(childHangConn) == "1" {
		// Dialled here, not at start-up. net/http stops counting a StateNew
		// connection against Shutdown once it is more than five seconds old,
		// so a connection opened before the wait would age out on a slow CI
		// run and Shutdown would succeed - leaving the test asserting nothing.
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			fmt.Println("dial:", err)
			os.Exit(5)
		}
		defer func() { _ = c.Close() }()

		// And wait for the accept, for the opposite reason: an unaccepted
		// connection is not one Shutdown waits for either.
		select {
		case <-accepted:
		case <-time.After(10 * time.Second):
			fmt.Println("the server never accepted the stalling connection")
			os.Exit(6)
		}

		// A connection that has sent nothing keeps Shutdown busy: net/http
		// only treats a StateNew connection as idle once it is more than five
		// seconds old. A short budget makes the timeout deterministic without
		// waiting out the real one.
		b.server = 300 * time.Millisecond
	}

	started := time.Now()
	serverErr, cleanupErr := gracefulShutdown(srv, quit, disarm, b)
	spent := time.Since(started)

	if serverErr != nil {
		// Deliberately not fatal, and deliberately not a bare return: the
		// point is that whatever follows still runs.
		fmt.Println("shutdown error:", serverErr)
	} else {
		fmt.Println(markerShutdown)
	}
	if cleanupErr != nil {
		fmt.Println("cleanup error:", cleanupErr)
	}
	fmt.Println(markerTook, spent.Nanoseconds())
	fmt.Println(markerExiting)
	_ = os.Stdout.Sync()
}

func startChild(t *testing.T, stuck bool, extraEnv ...string) (*exec.Cmd, chan string) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestSignalChild", "-test.v")
	cmd.Env = append(os.Environ(), childEnv+"=1")
	if stuck {
		cmd.Env = append(cmd.Env, childStuckEnv+"=1")
	}
	cmd.Env = append(cmd.Env, extraEnv...)
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	_ = w.Close()

	lines := make(chan string, 256)
	go func() {
		defer close(lines)
		buf := make([]byte, 4096)
		var acc strings.Builder
		for {
			n, err := r.Read(buf)
			if n > 0 {
				acc.Write(buf[:n])
				for {
					s := acc.String()
					i := strings.IndexByte(s, '\n')
					if i < 0 {
						break
					}
					lines <- s[:i]
					acc.Reset()
					acc.WriteString(s[i+1:])
				}
			}
			if err != nil {
				if acc.Len() > 0 {
					lines <- acc.String()
				}
				return
			}
		}
	}()

	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		_ = r.Close()
	})
	return cmd, lines
}

// await drains lines until one contains want, or the deadline passes. It
// returns everything it saw, so a failure says what the child actually did,
// and the matching line, so a marker can carry a value.
func await(t *testing.T, lines chan string, want string, d time.Duration) ([]string, string) {
	t.Helper()
	var seen []string
	deadline := time.After(d)
	for {
		select {
		case l, ok := <-lines:
			if !ok {
				t.Fatalf("child output ended before %q; saw:\n%s", want, strings.Join(seen, "\n"))
			}
			seen = append(seen, l)
			if strings.Contains(l, want) {
				return seen, l
			}
		case <-deadline:
			t.Fatalf("timed out waiting for %q; saw:\n%s", want, strings.Join(seen, "\n"))
		}
	}
}

// childAddr waits for the address the child is listening on.
func childAddr(t *testing.T, lines chan string) string {
	t.Helper()
	_, line := await(t, lines, markerAddr, 30*time.Second)
	fields := strings.Fields(line)
	return fields[len(fields)-1]
}

// took reads the nanoseconds gracefulShutdown spent, as the child measured
// them. Measured inside the child on purpose: the parent's own clock includes
// process scheduling, which is the noise the tightest assertion here cannot
// afford.
func took(t *testing.T, lines chan string, d time.Duration) time.Duration {
	t.Helper()
	_, line := await(t, lines, markerTook, d)
	fields := strings.Fields(line)
	ns, err := strconv.ParseInt(fields[len(fields)-1], 10, 64)
	if err != nil {
		t.Fatalf("unreadable %s line %q: %v", markerTook, line, err)
	}
	return time.Duration(ns)
}

// sample is one answer, or the refusal that replaced it.
type sample struct {
	at   time.Time
	path string
	// status is zero when the connection could not be made at all, which is
	// what a closed listener looks like from outside.
	status   int
	draining bool
	// willClose is what the server answered about the connection: the header
	// it sends is Connection: close, which the transport consumes and reports
	// here rather than leaving in Response.Header.
	willClose bool
}

// probe asks once, on a connection of its own.
//
// A new transport per request, because a connection opened before the signal
// can still be served after the listener is closed: reusing one would let this
// test pass against a shutdown that had already broken the listener. Keep-alive
// is left enabled so the server's own Connection: close is observable - a
// client that asked for close would get that header back either way, and the
// assertion would prove nothing.
func probe(addr, path string) sample {
	tr := &http.Transport{}
	defer tr.CloseIdleConnections()
	c := &http.Client{Transport: tr, Timeout: 3 * time.Second}

	s := sample{at: time.Now(), path: path}
	resp, err := c.Get("http://" + addr + path)
	if err != nil {
		return s
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	s.status = resp.StatusCode
	s.willClose = resp.Close
	s.draining = strings.Contains(string(body), `"status":"draining"`)
	return s
}

// watcher polls the child until it stops accepting connections, keeping every
// answer.
type watcher struct {
	mu      sync.Mutex
	samples []sample
	done    chan struct{}
}

func watch(addr string, paths ...string) *watcher {
	w := &watcher{done: make(chan struct{})}
	go func() {
		defer close(w.done)
		for {
			refused := false
			for _, p := range paths {
				s := probe(addr, p)
				w.mu.Lock()
				w.samples = append(w.samples, s)
				w.mu.Unlock()
				if s.status == 0 {
					refused = true
				}
			}
			if refused {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	return w
}

// sawDraining reports whether /ready has answered "draining" yet.
func (w *watcher) sawDraining() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, s := range w.samples {
		if s.path == readyPath && s.draining {
			return true
		}
	}
	return false
}

func (w *watcher) wait(t *testing.T, d time.Duration) []sample {
	t.Helper()
	select {
	case <-w.done:
	case <-time.After(d):
		t.Fatal("the child never stopped accepting connections")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.samples
}

func describe(samples []sample) string {
	var b strings.Builder
	for _, s := range samples {
		fmt.Fprintf(&b, "  %s %s -> %d draining=%v willClose=%v\n",
			s.at.Format("15:04:05.000"), s.path, s.status, s.draining, s.willClose)
	}
	return b.String()
}

// Acceptance 19. Registering only os.Interrupt meant SIGTERM - the signal
// `docker stop`, Kubernetes and systemd all send - terminated the process
// before any of the shutdown path ran. Both must now reach it.
func TestBothSignalsRunTheShutdownPath(t *testing.T) {
	for _, tc := range []struct {
		name string
		sig  syscall.Signal
	}{
		{"SIGINT", syscall.SIGINT},
		{"SIGTERM", syscall.SIGTERM},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd, lines := startChild(t, false)
			await(t, lines, markerReady, 30*time.Second)

			if err := cmd.Process.Signal(tc.sig); err != nil {
				t.Fatalf("signal: %v", err)
			}

			await(t, lines, markerSignal, 10*time.Second)
			await(t, lines, markerShutdown, 10*time.Second)
			await(t, lines, markerExiting, 10*time.Second)

			if err := cmd.Wait(); err != nil {
				t.Fatalf("child exited with %v, want a clean exit", err)
			}
		})
	}
}

// Acceptance 20. quit is a buffered channel and signal.Notify stays armed, so
// without restoring the disposition a second signal only refills the buffer:
// once SIGTERM is registered, a shutdown that hangs could not be interrupted by
// anything short of SIGKILL.
//
// The hang is now a cleanup callback that never returns, which is where a
// shutdown actually hangs, and it is reached through gracefulShutdown - so this
// also pins where the disposition is restored. Restore it before the drain
// window and the window itself becomes the interruptible part; restore it never
// and this test hangs.
func TestASecondSignalStillKillsAStuckShutdown(t *testing.T) {
	cmd, lines := startChild(t, true)
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("first signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	// The child is now on its way into a cleanup that will not finish on its
	// own. Signalled repeatedly rather than once: the marker is printed just
	// before gracefulShutdown is entered, and the disposition is not restored
	// until the drain window is over - zero seconds here, but not zero
	// instructions - so a single signal sent immediately after the marker can
	// still land in the buffered channel and be dropped. Which of them does
	// the killing is not the assertion; that one of them can is.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	retry := time.NewTicker(200 * time.Millisecond)
	defer retry.Stop()
	deadline := time.After(15 * time.Second)
	for {
		select {
		case err := <-done:
			if err == nil {
				t.Fatal("child exited cleanly; it was supposed to be killed by the second signal")
			}
			return
		case <-retry.C:
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				t.Fatalf("second signal: %v", err)
			}
		case <-deadline:
			t.Fatal("the second signal did not kill a stuck shutdown - the escape hatch is gone")
		}
	}
}

// Acceptance 21. srv.Shutdown reports an error exactly when connections were
// still in flight, and the old code answered that with log.Fatal - an
// unconditional os.Exit(1). Everything after it, which is where the cleanup
// hooks will hang, never ran. A failed Shutdown must not end the process.
func TestShutdownTimeoutDoesNotStopWhatFollows(t *testing.T) {
	cmd, lines := startChild(t, false, childHangConn+"=1")
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	seen, _ := await(t, lines, markerExiting, 20*time.Second)

	var timedOut bool
	for _, l := range seen {
		if strings.Contains(l, "shutdown error:") {
			timedOut = true
		}
	}
	if !timedOut {
		t.Fatalf("Shutdown did not time out, so this test proves nothing; saw:\n%s",
			strings.Join(seen, "\n"))
	}
	var cleaned bool
	for _, l := range seen {
		if strings.Contains(l, markerCleanup) {
			cleaned = true
		}
	}
	if !cleaned {
		t.Fatalf("the BeforeExit callback did not run after a failed Shutdown; saw:\n%s",
			strings.Join(seen, "\n"))
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with %v after a failed Shutdown, want a clean exit", err)
	}
}

// A callback that outlasts its budget must not take the process with it, and
// must not be waited for: RunShutdown reports the deadline and returns, the
// callback carries on, and the process still exits cleanly. This is the half of
// the contract that is easy to get backwards - the context bounds the wait, not
// the work, because Go cannot cancel a function that does not check for it.
func TestACleanupThatOutlastsItsBudgetIsAbandonedNotAwaited(t *testing.T) {
	cmd, lines := startChild(t, false, childSlowCleanup+"=1")
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	// The budget is 300ms and the callback sleeps two seconds. If RunShutdown
	// waited for it, this marker would not arrive for two seconds; the one
	// second here is what makes "abandoned, not awaited" the thing asserted.
	seen, _ := await(t, lines, markerExiting, 1*time.Second)

	var reported bool
	for _, l := range seen {
		if strings.Contains(l, "cleanup error:") {
			reported = true
		}
		if strings.Contains(l, markerCleanup) {
			t.Fatalf("the slow callback finished before the process moved on, so nothing was abandoned; saw:\n%s",
				strings.Join(seen, "\n"))
		}
	}
	if !reported {
		t.Fatalf("RunShutdown returned no error for a callback that outlasted the budget; saw:\n%s",
			strings.Join(seen, "\n"))
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with %v, want a clean exit despite the abandoned callback", err)
	}
}

// The core acceptance: with a drain window configured, something outside the
// process can observe that this instance is draining, on a connection it opens
// after the signal, and can still be served while it does.
//
// Two windows rather than one. A single value proves only that something takes
// that long, which a hard-coded sleep anywhere in the sequence would satisfy;
// two say the wait is the configured one.
//
// What each answer is for:
//
//   - /ready reporting "draining" is the window being observable at all. The
//     status code alone would not say it: with no database configured the
//     probe's own checks fail and 503 is also the answer before the signal.
//   - The server refusing to keep those connections alive is the window being
//     useful. It keeps them alive until Shutdown sets shuttingDown(), so
//     without switching keep-alive off here a balancer's pool would sit
//     untouched for the whole window and be cut at the end of it anyway. The
//     header saying so is Connection: close; the transport consumes it and
//     reports it as Response.Close, which is what a sample records.
//   - /health staying 200 is the window not asking to be restarted, and the
//     ordinary route staying 200 is the window not refusing work. Draining is
//     "stop sending me new work", not "reject what arrives".
func TestTheDrainWindowIsObservableWhileStillServing(t *testing.T) {
	for _, drain := range []time.Duration{300 * time.Millisecond, 1200 * time.Millisecond} {
		t.Run(drain.String(), func(t *testing.T) {
			cmd, lines := startChild(t, false,
				fmt.Sprintf("%s=%d", childDrainMS, drain.Milliseconds()))
			addr := childAddr(t, lines)
			await(t, lines, markerReady, 30*time.Second)

			w := watch(addr, readyPath, healthPath, pingPath)
			// Long enough for a round of answers from a server that is not yet
			// draining, which is what the keep-alive assertion below compares
			// against.
			time.Sleep(150 * time.Millisecond)

			signalAt := time.Now()
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				t.Fatalf("signal: %v", err)
			}

			samples := w.wait(t, drain+30*time.Second)
			spent := took(t, lines, 10*time.Second)
			await(t, lines, markerExiting, 10*time.Second)

			if spent < drain {
				t.Errorf("the shutdown took %s, want at least the %s window", spent, drain)
			}

			var refusedAt = -1
			for i, s := range samples {
				if s.status == 0 {
					refusedAt = i
					break
				}
			}
			if refusedAt < 0 {
				t.Fatalf("the child never stopped accepting; saw:\n%s", describe(samples))
			}

			var keptAliveBefore, drainingInside, closedInside bool
			for _, s := range samples[:refusedAt] {
				switch s.path {
				case readyPath:
					if s.at.Before(signalAt) && !s.draining && !s.willClose {
						keptAliveBefore = true
					}
					if s.at.After(signalAt) && s.draining {
						drainingInside = true
						if s.willClose {
							closedInside = true
						}
					}
				case healthPath, pingPath:
					if s.status != http.StatusOK {
						t.Errorf("%s answered %d before the listener closed, want 200;\n%s",
							s.path, s.status, describe(samples))
					}
				}
			}

			if !keptAliveBefore {
				t.Fatalf("no answer before the signal kept the connection alive, so the header assertion below proves nothing;\n%s",
					describe(samples))
			}
			if !drainingInside {
				t.Errorf("no answer inside the window reported draining; the flip and the closed listener were not far enough apart to observe;\n%s",
					describe(samples))
			}
			if !closedInside {
				t.Errorf("answers inside the window still kept the connection alive, so a pooled connection survives the whole window and is cut at the end of it anyway;\n%s",
					describe(samples))
			}

			if err := cmd.Wait(); err != nil {
				t.Fatalf("child exited with %v, want a clean exit", err)
			}
		})
	}
}

// The default has to be no window at all: a process that configures no
// extend.shutdown section must shut down the way it did before the section
// existed.
//
// Asserted as a sequence rather than as a duration. How long a shutdown takes
// is decided by how much the cleanup callbacks have to do, so "as fast as
// before" is not falsifiable; "nothing was inserted between the signal and the
// listener closing" is.
func TestAnUnconfiguredShutdownAddsNoWindow(t *testing.T) {
	cmd, lines := startChild(t, false)
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	spent := took(t, lines, 10*time.Second)
	if spent > 100*time.Millisecond {
		t.Errorf("an unconfigured shutdown spent %s between the signal and exiting; "+
			"with no drain window and no cleanup callbacks it must be immediate", spent)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with %v, want a clean exit", err)
	}
}

// A second signal during the window ends it early rather than killing the
// process. Somebody sending another kill wants this over with sooner, and the
// answer to that is to stop draining - not to skip the cleanup, which is what
// the default disposition would do.
//
// This is the pair to TestASecondSignalStillKillsAStuckShutdown: the escape
// hatch has to be closed for the length of the window and open after it.
func TestASecondSignalEndsTheDrainWindowEarly(t *testing.T) {
	// Long enough that the shutdown cannot plausibly have taken this long on
	// its own, short enough that the test does not sit out the whole window
	// when the early exit is missing - it fails on the reported duration
	// instead of on a timeout, which says which of the two broke.
	const window = 10 * time.Second
	cmd, lines := startChild(t, false,
		fmt.Sprintf("%s=%d", childDrainMS, window.Milliseconds()))
	addr := childAddr(t, lines)
	await(t, lines, markerReady, 30*time.Second)

	w := watch(addr, readyPath)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("first signal: %v", err)
	}

	deadline := time.Now().Add(15 * time.Second)
	for !w.sawDraining() {
		if time.Now().After(deadline) {
			t.Fatal("the child never reported draining, so the second signal below would not land inside the window")
		}
		time.Sleep(20 * time.Millisecond)
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("second signal: %v", err)
	}

	spent := took(t, lines, window+20*time.Second)
	if spent >= window {
		t.Errorf("the window ran its full %s despite a second signal (%s); the signal was ignored", window, spent)
	}
	await(t, lines, markerExiting, 10*time.Second)

	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with %v; a second signal inside the window must end the window, not the process", err)
	}
}
