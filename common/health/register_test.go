package health

import (
	"context"
	"errors"
	"testing"
)

// isolate empties the registered checks and puts them back, so one test in
// this package cannot decide what the next one sees.
func isolate(t *testing.T) {
	t.Helper()
	extraMu.Lock()
	previous := extra
	extra = nil
	extraMu.Unlock()
	t.Cleanup(func() {
		extraMu.Lock()
		extra = previous
		extraMu.Unlock()
	})
}

func findCheck(checks []Check, name string) (Check, bool) {
	for _, c := range checks {
		if c.Name == name {
			return c, true
		}
	}
	return Check{}, false
}

// A registered check has to reach Ready's answer, or the host has wired
// something that never gets asked.
func TestARegisteredCheckIsAsked(t *testing.T) {
	isolate(t)
	Register("schema", func(context.Context) error { return errors.New("two behind") })

	checks := Ready(context.Background())
	c, ok := findCheck(checks, "schema")
	if !ok {
		t.Fatal("Ready did not ask the registered check")
	}
	if c.OK {
		t.Error("the check returned an error and was still reported OK")
	}
	if c.Err != "two behind" {
		t.Errorf("Err = %q, want the check's own message", c.Err)
	}
	if Healthy(checks) {
		t.Error("Healthy said yes while a registered check was failing")
	}
}

// The context Ready is given has to reach the check: it carries the probe's
// deadline, and a check that ignores it can hold the handler past it.
func TestTheRegisteredCheckIsGivenReadysContext(t *testing.T) {
	isolate(t)
	type key struct{}
	Register("ctx", func(ctx context.Context) error {
		if ctx.Value(key{}) != "carried" {
			return errors.New("the check was handed a different context")
		}
		return nil
	})

	checks := Ready(context.WithValue(context.Background(), key{}, "carried"))
	c, ok := findCheck(checks, "ctx")
	if !ok {
		t.Fatal("the registered check was not asked")
	}
	if !c.OK {
		t.Errorf("check failed: %s", c.Err)
	}
}

// A check that panics must not take the process down through the probe, the
// same guarantee the built-in checks have.
func TestARegisteredCheckThatPanicsFailsRatherThanCrashes(t *testing.T) {
	isolate(t)
	Register("boom", func(context.Context) error { panic("registry unreachable") })

	checks := Ready(context.Background())
	c, ok := findCheck(checks, "boom")
	if !ok {
		t.Fatal("the registered check was not asked")
	}
	if c.OK {
		t.Error("a panicking check was reported OK")
	}
}

func TestRegisteringTheSameNameTwicePanics(t *testing.T) {
	isolate(t)
	Register("dup", func(context.Context) error { return nil })

	defer func() {
		if recover() == nil {
			t.Error("registering a duplicate name did not panic; two checks under " +
				"one name make the failing one impossible to identify")
		}
	}()
	Register("dup", func(context.Context) error { return nil })
}

func TestRegisterRefusesAnEmptyNameOrNilCheck(t *testing.T) {
	isolate(t)
	for _, tc := range []struct {
		name string
		fn   func(context.Context) error
		why  string
	}{
		{"", func(context.Context) error { return nil }, "empty name"},
		{"nilfn", nil, "nil function"},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", tc.why)
				}
			}()
			Register(tc.name, tc.fn)
		}()
	}
}
