package api

import (
	"testing"
	"time"

	ext "go-admin/config"
)

// The seconds in the configuration and the durations the sequence waits on are
// two spellings of one budget.
func TestBudgetFromSeconds(t *testing.T) {
	got := budgetFrom(ext.ShutdownBudget{Server: 5, Cleanup: 3})
	want := budget{
		server:  5 * time.Second,
		cleanup: 3 * time.Second,
	}
	if got != want {
		t.Errorf("budgetFrom = %+v, want %+v", got, want)
	}
}

// The package variables and config.Default*Seconds have to say the same thing.
// They are the same default written twice - once as durations for the shutdown
// and once as seconds for the fallback - and a deployment that configures
// nothing is entitled to one answer, not two.
func TestDefaultBudgetIsTheConfiguredFallback(t *testing.T) {
	unconfigured, err := ext.Shutdown{}.Budget()
	if err != nil {
		t.Fatalf("the empty section did not resolve: %v", err)
	}
	if got, want := defaultBudget(), budgetFrom(unconfigured); got != want {
		t.Errorf("defaultBudget = %+v, want the unconfigured budget %+v", got, want)
	}
}
