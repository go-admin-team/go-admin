package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/gin-gonic/gin"

	"go-admin/config"
)

// serve builds a router with the limiter in front of a handler that always
// succeeds, so any non-200 comes from the limiter.
func serve(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Sentinel())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func get(t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
	return w
}

// TestSentinelRejectsWithTooManyRequests pins the status code. A rejected
// request used to answer 200 with the failure only in the body, so every layer
// that reads the status line - load balancers, metrics, client retry, load
// tests - counted it as served.
func TestSentinelRejectsWithTooManyRequests(t *testing.T) {
	one := 1.0
	config.ExtConfig.RateLimit = config.RateLimit{InboundQPS: &one}
	t.Cleanup(func() {
		config.ExtConfig.RateLimit = config.RateLimit{}
		_ = system.ClearRules()
	})

	r := serve(t)

	var rejected *httptest.ResponseRecorder
	for i := 0; i < 20; i++ {
		if w := get(t, r); w.Code != http.StatusOK {
			rejected = w
			break
		}
	}
	if rejected == nil {
		t.Fatal("a limit of 1 req/s let 20 requests through; the limiter is not engaged")
	}
	if rejected.Code != http.StatusTooManyRequests {
		t.Errorf("rejected with %d, want %d", rejected.Code, http.StatusTooManyRequests)
	}

	// The body's code must agree with the status line; they disagreed before.
	var body struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(rejected.Body.Bytes(), &body); err != nil {
		t.Fatalf("rejection body is not json: %v", err)
	}
	if body.Code != http.StatusTooManyRequests {
		t.Errorf("body code = %d, want %d", body.Code, http.StatusTooManyRequests)
	}
	if body.Msg == "" {
		t.Error("rejection carries no message")
	}
}

// TestSentinelDisabledByZero covers the escape hatch: a deployment behind its
// own gateway has no use for a second limiter.
//
// It asserts on the loaded rules rather than on traffic. Sentinel measures QPS
// over a sliding window, so a burst issued inside one bucket is not counted
// before the bucket closes - a few hundred requests sail past a threshold of
// 200 in a test, and "no request was rejected" would pass whether or not the
// limiter is disabled. Whether a rule was installed at all does not depend on
// timing.
func TestSentinelDisabledByZero(t *testing.T) {
	if err := system.ClearRules(); err != nil {
		t.Fatal(err)
	}
	zero := 0.0
	config.ExtConfig.RateLimit = config.RateLimit{InboundQPS: &zero}
	t.Cleanup(func() {
		config.ExtConfig.RateLimit = config.RateLimit{}
		_ = system.ClearRules()
	})

	r := serve(t)
	if rules := system.GetRules(); len(rules) != 0 {
		t.Errorf("limiter disabled but %d rule(s) were loaded: %+v", len(rules), rules)
	}

	for i := 0; i < 500; i++ {
		if w := get(t, r); w.Code != http.StatusOK {
			t.Fatalf("request %d got %d with the limiter disabled", i, w.Code)
		}
	}
}
