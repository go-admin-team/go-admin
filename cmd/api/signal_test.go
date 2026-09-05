package api

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The signal path cannot be exercised in-process: delivering a signal to the
// test binary would race with the test framework, and the disposition changes
// are global. So the test re-executes itself as a child, and the child runs the
// same armStopSignals / shutdownServer the server does.
//
// The child deliberately serves an empty http.Server rather than the real one:
// this repository's CI has no database (.github/workflows/go.yml runs neither
// MySQL nor a sqlite-tagged build), and none of what is under test needs one.
const (
	childEnv       = "GO_ADMIN_SIGNAL_CHILD"
	childStuckEnv  = "GO_ADMIN_SIGNAL_CHILD_STUCK"
	childHangConn  = "GO_ADMIN_SIGNAL_CHILD_HANGCONN"
	markerReady    = "CHILD-READY"
	markerSignal   = "CHILD-SIGNAL"
	markerShutdown = "CHILD-SHUTDOWN-OK"
	markerExiting  = "CHILD-EXITING"
)

// TestSignalChild is the child process. It is skipped in a normal run.
func TestSignalChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("child process entry point")
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		os.Exit(3)
	}
	srv := &http.Server{Handler: http.NewServeMux()}
	go func() { _ = srv.Serve(ln) }()

	// Arm before announcing readiness. Doing it the other way round leaves a
	// window in which the parent's signal reaches the default handler and
	// kills the child before any of this runs - which is exactly the failure
	// this whole change is about, so the test must not reproduce it by
	// accident.
	quit, disarm := armStopSignals()

	fmt.Println(markerReady)
	os.Stdout.Sync()

	sig := <-quit
	disarm()
	fmt.Println(markerSignal, sig)
	os.Stdout.Sync()

	if os.Getenv(childStuckEnv) == "1" {
		// Stand in for a cleanup hook that never finishes. The point of
		// restoring the signal disposition is that a second signal still
		// reaches the default handler and kills this.
		time.Sleep(2 * time.Minute)
	}

	timeout := shutdownTimeout
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

		// A connection that has sent nothing keeps Shutdown busy: net/http
		// only treats a StateNew connection as idle once it is more than five
		// seconds old. A short budget makes the timeout deterministic without
		// waiting out the real one.
		timeout = 300 * time.Millisecond
	}

	if err := shutdownServer(srv, timeout); err != nil {
		// Deliberately not fatal, and deliberately not a bare return: the
		// point is that whatever follows still runs.
		fmt.Println("shutdown error:", err)
	} else {
		fmt.Println(markerShutdown)
	}
	fmt.Println(markerExiting)
	os.Stdout.Sync()
}

func startChild(t *testing.T, stuck bool, extraEnv ...string) (*exec.Cmd, *os.File, chan string) {
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

	lines := make(chan string, 64)
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
	return cmd, r, lines
}

// await drains lines until one contains want, or the deadline passes. It
// returns everything it saw, so a failure says what the child actually did.
func await(t *testing.T, lines chan string, want string, d time.Duration) []string {
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
				return seen
			}
		case <-deadline:
			t.Fatalf("timed out waiting for %q; saw:\n%s", want, strings.Join(seen, "\n"))
		}
	}
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
			cmd, _, lines := startChild(t, false)
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
func TestASecondSignalStillKillsAStuckShutdown(t *testing.T) {
	cmd, _, lines := startChild(t, true)
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("first signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	// The child is now inside a cleanup that will not finish on its own.
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("second signal: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("child exited cleanly; it was supposed to be killed by the second signal")
		}
	case <-time.After(15 * time.Second):
		t.Fatal("the second signal did not kill a stuck shutdown - the escape hatch is gone")
	}
}

// Acceptance 21. srv.Shutdown reports an error exactly when connections were
// still in flight, and the old code answered that with log.Fatal - an
// unconditional os.Exit(1). Everything after it, which is where the cleanup
// hooks will hang, never ran. A failed Shutdown must not end the process.
func TestShutdownTimeoutDoesNotStopWhatFollows(t *testing.T) {
	cmd, _, lines := startChild(t, false, childHangConn+"=1")
	await(t, lines, markerReady, 30*time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("signal: %v", err)
	}
	await(t, lines, markerSignal, 10*time.Second)

	seen := await(t, lines, markerExiting, 20*time.Second)

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
	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with %v after a failed Shutdown, want a clean exit", err)
	}
}
