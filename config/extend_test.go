package config

import (
	"strings"
	"testing"
)

func TestObjectStoreConfigured(t *testing.T) {
	if (ObjectStore{}).Configured() {
		t.Fatal("empty store reported as configured")
	}
	full := ObjectStore{Endpoint: "e", AccessKeyID: "a", AccessKeySecret: "s", BucketName: "b"}
	if !full.Configured() {
		t.Fatal("complete store reported as unconfigured")
	}
	if (ObjectStore{Endpoint: "e", AccessKeyID: "a"}).Configured() {
		t.Fatal("partial store reported as configured")
	}
}

func TestRateLimitThreshold(t *testing.T) {
	// Absent is the case an existing settings.yml hits after an upgrade: it has
	// no ratelimit section, and must keep the limit it always had.
	if got := (RateLimit{}).Threshold(); got != DefaultInboundQPS {
		t.Errorf("unconfigured limit = %v, want the default %v", got, DefaultInboundQPS)
	}

	zero := 0.0
	if got := (RateLimit{InboundQPS: &zero}).Threshold(); got != 0 {
		t.Errorf("explicit zero = %v, want 0 so the limiter can be turned off", got)
	}

	custom := 1500.0
	if got := (RateLimit{InboundQPS: &custom}).Threshold(); got != custom {
		t.Errorf("configured limit = %v, want %v", got, custom)
	}
}

func ptr(v int) *int { return &v }

// The zero-value rule is the one the section would otherwise need a paragraph
// of documentation to survive: nil takes the default, a number that was
// written down is spent literally. A `server: 0` that quietly became five
// seconds would be configuration accepted and not applied.
func TestShutdownBudgetFallbacks(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   Shutdown
		want ShutdownBudget
	}{
		{
			// What an existing settings.yml hits after an upgrade: no
			// extend.shutdown section at all, and therefore the shutdown it
			// already had.
			name: "nothing configured",
			in:   Shutdown{},
			want: ShutdownBudget{Server: 5, Cleanup: 3},
		},
		{
			name: "both configured",
			in:   Shutdown{Server: ptr(8), Cleanup: ptr(4)},
			want: ShutdownBudget{Server: 8, Cleanup: 4},
		},
		{
			// The case a plain int could not express: do not wait for
			// in-flight requests, which is a reasonable thing to ask when the
			// grace period is very short.
			name: "explicit zeros are spent, not replaced",
			in:   Shutdown{Server: ptr(0), Cleanup: ptr(0)},
			want: ShutdownBudget{Server: 0, Cleanup: 0},
		},
		{
			name: "one field configured, the rest default",
			in:   Shutdown{Cleanup: ptr(15)},
			want: ShutdownBudget{Server: 5, Cleanup: 15},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.in.Budget()
			if err != nil {
				t.Fatalf("Budget() = %v", err)
			}
			if got != tc.want {
				t.Errorf("Budget() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

// A negative is refused, not corrected. Turning it into zero would be the
// failure this section exists to remove - written down, accepted, and not what
// happens - and there is no reading of a negative wait to honour.
//
// The last row is what makes the others mean anything: an implementation that
// refused every value would pass them all.
func TestShutdownBudgetRefusesNegativeSeconds(t *testing.T) {
	for _, tc := range []struct {
		name    string
		in      Shutdown
		wantErr bool
	}{
		{name: "negative server", in: Shutdown{Server: ptr(-1)}, wantErr: true},
		{name: "negative cleanup", in: Shutdown{Cleanup: ptr(-1)}, wantErr: true},
		{name: "explicit zeros are not negative", in: Shutdown{Server: ptr(0), Cleanup: ptr(0)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.in.Budget()
			if tc.wantErr && err == nil {
				t.Fatal("Budget() accepted a negative number of seconds")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Budget() = %v, want the zeros taken literally", err)
			}
		})
	}
}

// The message has to name every field that is wrong, not the first one: a
// caller who fixes one and gets the same error back learns to distrust it.
func TestShutdownBudgetNamesEveryNegativeField(t *testing.T) {
	_, err := Shutdown{Server: ptr(-30), Cleanup: ptr(-3)}.Budget()
	if err == nil {
		t.Fatal("Budget() accepted two negative values")
	}
	for _, name := range []string{"server", "cleanup"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("%q does not name %s", err, name)
		}
	}
}
