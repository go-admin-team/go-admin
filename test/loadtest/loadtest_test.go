// Package loadtest measures what one go-admin process sustains over HTTP.
//
// It is skipped unless GOADMIN_BENCH_ADDR points at a running server, so
// `go test ./...` is unaffected. Start a server and run:
//
//	GOADMIN_BENCH_ADDR=http://127.0.0.1:8000 go test ./test/loadtest/ -v -run TestLoadProfile
//
// Unlike a Go benchmark this reports latency percentiles, which is what
// capacity planning needs: an average hides the tail that users actually feel.
//
// Two caveats when reading the numbers. The load generator runs on the same
// machine as the server unless GOADMIN_BENCH_ADDR is remote, so both compete
// for the same cores - a split deployment measures higher. And the figures
// describe the configured backend: sqlite and MySQL differ by more than the
// framework does.
package loadtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	addrEnv  = "GOADMIN_BENCH_ADDR"
	tokenEnv = "GOADMIN_BENCH_TOKEN"
	userEnv  = "GOADMIN_BENCH_USER"
	passEnv  = "GOADMIN_BENCH_PASS"

	// Each concurrency level runs for this long. Long enough to get past
	// connection setup and let the scheduler settle, short enough that the
	// whole sweep stays interactive.
	levelDuration = 3 * time.Second
)

// concurrencyLevels sweeps from a single client to well past core count, so
// the point where added concurrency stops buying throughput is visible rather
// than assumed. Peak throughput and peak concurrency are not the same number:
// past the peak a server takes more work than it can finish and both
// throughput and latency get worse, so the sweep has to bracket the turn
// rather than stop at the top.
//
// GOADMIN_BENCH_LEVELS overrides it, comma separated.
var concurrencyLevels = parseLevels(os.Getenv("GOADMIN_BENCH_LEVELS"), []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512})

func parseLevels(spec string, fallback []int) []int {
	if spec == "" {
		return fallback
	}
	out := make([]int, 0, 8)
	for _, f := range strings.Split(spec, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil || n <= 0 {
			continue
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

func addr(t testing.TB) string {
	t.Helper()
	a := os.Getenv(addrEnv)
	if a == "" {
		t.Skipf("%s not set; skipping load test", addrEnv)
	}
	return a
}

// newClient returns a client whose pool is large enough that the generator
// does not become the bottleneck it is trying to measure.
func newClient(maxConns int) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        maxConns * 2,
			MaxIdleConnsPerHost: maxConns * 2,
			MaxConnsPerHost:     maxConns * 2,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
		},
	}
}

// result is one completed request.
type result struct {
	latency time.Duration
	err     bool
	status  int
}

// report is the summary of one concurrency level.
type report struct {
	concurrency        int
	total              int64
	failed             int64
	elapsed            time.Duration
	p50, p95, p99, max time.Duration
	statuses           map[int]int64
}

func (r report) qps() float64 {
	if r.elapsed == 0 {
		return 0
	}
	return float64(r.total) / r.elapsed.Seconds()
}

func (r report) String() string {
	codes := make([]int, 0, len(r.statuses))
	for c := range r.statuses {
		codes = append(codes, c)
	}
	sort.Ints(codes)
	dist := make([]string, 0, len(codes))
	for _, c := range codes {
		dist = append(dist, fmt.Sprintf("%d:%d", c, r.statuses[c]))
	}
	return fmt.Sprintf("c=%-4d %9.0f req/s  p50=%-9s p95=%-9s p99=%-9s max=%-9s failed=%-7d %s",
		r.concurrency, r.qps(),
		r.p50.Round(time.Microsecond), r.p95.Round(time.Microsecond),
		r.p99.Round(time.Microsecond), r.max.Round(time.Microsecond), r.failed,
		strings.Join(dist, " "))
}

// drive runs `concurrency` workers against req for levelDuration and collects
// every latency. Bodies are drained and closed - skipping that silently caps
// throughput at the point connections stop being reused.
func drive(t testing.TB, concurrency int, want int, mk func() *http.Request) report {
	t.Helper()

	client := newClient(concurrency)
	defer client.CloseIdleConnections()

	ctx, cancel := context.WithTimeout(context.Background(), levelDuration)
	defer cancel()

	var (
		mu       sync.Mutex
		samples  []time.Duration
		statuses = map[int]int64{}
		failed   atomic.Int64
		total    atomic.Int64
		wg       sync.WaitGroup
	)

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := make([]time.Duration, 0, 1024)
			localStatus := map[int]int64{}
			for ctx.Err() == nil {
				req := mk()
				t0 := time.Now()
				resp, err := client.Do(req.WithContext(ctx))
				d := time.Since(t0)
				if err != nil {
					if ctx.Err() != nil {
						break
					}
					failed.Add(1)
					total.Add(1)
					continue
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				localStatus[resp.StatusCode]++
				if resp.StatusCode != want {
					failed.Add(1)
				}
				total.Add(1)
				local = append(local, d)
			}
			mu.Lock()
			samples = append(samples, local...)
			for code, n := range localStatus {
				statuses[code] += n
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	r := report{
		concurrency: concurrency,
		total:       total.Load(),
		failed:      failed.Load(),
		elapsed:     elapsed,
		statuses:    statuses,
	}
	if n := len(samples); n > 0 {
		r.p50 = samples[n*50/100]
		r.p95 = samples[min(n*95/100, n-1)]
		r.p99 = samples[min(n*99/100, n-1)]
		r.max = samples[n-1]
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// login obtains a token. GOADMIN_BENCH_TOKEN short-circuits it, which is how a
// server in prod mode is reached - there the login endpoint demands a captcha.
func login(t testing.TB, base string) string {
	t.Helper()
	if tok := os.Getenv(tokenEnv); tok != "" {
		return tok
	}

	user, pass := os.Getenv(userEnv), os.Getenv(passEnv)
	if user == "" {
		user, pass = "admin", "123456"
	}

	body, _ := json.Marshal(map[string]string{
		"username": user,
		"password": pass,
		"code":     "0",
		"uuid":     "0",
	})
	resp, err := http.Post(base+"/api/v1/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login returned %d: %s\n(a server in prod mode requires a captcha; set %s instead)",
			resp.StatusCode, raw, tokenEnv)
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.Token == "" {
		t.Fatalf("no token in login response: %s", raw)
	}
	return out.Token
}

// TestLoadProfile sweeps concurrency against three endpoints chosen for what
// they isolate:
//
//   - captcha: no auth, no business query. The routing and image-generation
//     floor.
//   - dept list: the full authenticated path - JWT parse, casbin check, data
//     permission scope, database read. This is what a real page costs.
//   - login: bcrypt. Deliberately slow, and the one endpoint whose ceiling is
//     set by design rather than by the framework.
func TestLoadProfile(t *testing.T) {
	base := addr(t)
	token := login(t, base)

	cases := []struct {
		name string
		want int
		mk   func() *http.Request
	}{
		{
			// The control. An unrouted path exercises the HTTP stack, gin's
			// tree lookup and nothing else, so it bounds every other row here.
			// When a business endpoint reaches this number, the measurement has
			// stopped describing the endpoint and started describing the
			// transport - or the load generator, when both share a machine.
			name: "404 (http+routing floor)",
			want: 404,
			mk: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, base+"/api/v1/__no_such_route__", nil)
				return req
			},
		},
		{
			// The framework on its own: global middleware chain, route lookup,
			// and a handler that only sets a status. No database, no cache.
			// Against the 404 row this isolates what the chain costs; against
			// the rows below it, what the business path adds.
			//
			// Numbers from any endpoint that touches a database describe the
			// database, the driver and the pool as much as the framework - the
			// MySQL sweeps here moved from collapsing at c=64 to 19k req/s at
			// c=512 on a pool setting alone, with the framework untouched.
			name: "health (framework only, no db)",
			want: 200,
			mk: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, base+"/api/v1/health", nil)
				return req
			},
		},
		{
			name: "captcha (no auth)",
			want: 200,
			mk: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, base+"/api/v1/captcha", nil)
				return req
			},
		},
		{
			name: "dept list (jwt+casbin+db)",
			want: 200,
			mk: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, base+"/api/v1/dept", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, c := range concurrencyLevels {
				t.Log(drive(t, c, tc.want, tc.mk))
			}
		})
	}
}

// TestLoginThroughput is separated because bcrypt saturates the CPU: running
// it alongside the others would distort them. It also writes a login-log row
// per attempt when logger.enableddb is on, so the number moves with that
// setting.
func TestLoginThroughput(t *testing.T) {
	base := addr(t)
	if os.Getenv(tokenEnv) != "" {
		t.Skip("token supplied; login endpoint presumably needs a captcha")
	}

	user, pass := os.Getenv(userEnv), os.Getenv(passEnv)
	if user == "" {
		user, pass = "admin", "123456"
	}
	body, _ := json.Marshal(map[string]string{
		"username": user, "password": pass, "code": "0", "uuid": "0",
	})

	mk := func() *http.Request {
		req, _ := http.NewRequest(http.MethodPost, base+"/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	for _, c := range []int{1, 4, 8, 16, 32, 64, 128} {
		t.Log(drive(t, c, http.StatusOK, mk))
	}
}
