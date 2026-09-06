package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	otherrouter "go-admin/app/other/router"
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

// The rate limiter must not answer for the probes.
//
// It is installed on the engine, so without this the probes are limited like
// any other route and answer 429 above the threshold. A liveness probe that
// collects 429s fails its threshold and the container is restarted - taking
// capacity out of a deployment that is already short of it and pushing the
// rest closer to the threshold. The limiter working exactly as designed is
// what would cause it.
//
// The stand-in rejects everything rather than being a real limiter: what is
// under test is which requests reach it, and a real one would need the traffic
// to cross a threshold before it said anything.
func TestTheProbesSkipTheRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var reached []string
	r := gin.New()
	r.Use(exemptProbes(func(c *gin.Context) {
		reached = append(reached, c.FullPath())
		c.AbortWithStatus(http.StatusTooManyRequests)
	}))
	v1 := r.Group(otherrouter.APIPrefix)
	otherrouter.RegisterMonitorRouter(v1)
	v1.GET("/business", func(c *gin.Context) { c.Status(http.StatusOK) })

	for _, tc := range []struct {
		path    string
		limited bool
	}{
		{otherrouter.APIPrefix + otherrouter.HealthPath, false},
		{otherrouter.APIPrefix + otherrouter.ReadyPath, false},
		// Not a probe, and deliberately not exempt: the exemption is for the
		// two routes an orchestrator acts on, not for everything under
		// /api/v1 that happens to be unauthenticated.
		{otherrouter.APIPrefix + "/metrics", true},
		{otherrouter.APIPrefix + "/business", true},
	} {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if tc.limited {
				if w.Code != http.StatusTooManyRequests {
					t.Errorf("answered %d, want the middleware's 429 - it was skipped for a route that is not a probe", w.Code)
				}
				return
			}
			if w.Code == http.StatusTooManyRequests {
				t.Errorf("answered 429; a probe that can be rate-limited gets the container restarted under load")
			}
		})
	}

	// Said separately, because a probe could also answer 429 by itself: what
	// has to be true is that the middleware never saw the request.
	for _, p := range reached {
		if probePaths[p] {
			t.Errorf("the middleware ran for %s", p)
		}
	}
}
