package middleware

import "testing"

// excluded is excludedFromCasbin with the error dropped: these tests are about
// the answer and its cost, and CasbinExclude has no malformed entry to report.
func excluded(t testing.TB, path, method string) bool {
	t.Helper()
	ok, err := excludedFromCasbin(method, path)
	if err != nil {
		t.Fatalf("CasbinExclude holds a pattern that will not compile: %s", err)
	}
	return ok
}

// TestCasbinExcludeScanMatches pins the behaviour the scan has to keep: an
// excluded route is recognised, a protected one is not, and the method has to
// agree.
func TestCasbinExcludeScanMatches(t *testing.T) {
	cases := []struct {
		path, method string
		want         bool
	}{
		{"/api/v1/health", "GET", true},
		{"/api/v1/login", "POST", true},
		{"/api/v1/roleMenuTreeselect/12", "GET", true},
		{"/api/v1/dept", "GET", false},
		{"/api/v1/sys-user", "GET", false},
		// Same path, wrong method: sys-user is excluded for PUT only.
		{"/api/v1/sys-user", "PUT", true},
		{"/api/v1/health", "POST", false},
	}
	for _, c := range cases {
		if got := excluded(t, c.path, c.method); got != c.want {
			t.Errorf("excludedFromCasbin(%s %s) = %v, want %v", c.method, c.path, got, c.want)
		}
	}
}

// TestCasbinExcludeScanAllocationBudget is what keeps the scan cheap.
//
// The list is walked per request with a pattern match per entry, and
// casbin's util.KeyMatch2 compiles a regexp on every call - the whole scan
// cost about 2,566 allocations that way. Going back to it fails this test.
//
// Allocation counts are deterministic across machines; wall-clock is not.
func TestCasbinExcludeScanAllocationBudget(t *testing.T) {
	// A protected route, so the scan runs to the end without an early match -
	// the case every authenticated business request hits.
	const path, method = "/api/v1/dept", "GET"

	if excluded(t, path, method) {
		t.Fatalf("setup failed: %s is in the exclusion list", path)
	}

	// The budget covers the GET entries that carry a path parameter, which
	// still need a match. Measured at 0 for the cached matcher; the headroom
	// is for entries being added to the list.
	const budget = 64

	got := testing.AllocsPerRun(100, func() {
		_, _ = excludedFromCasbin(method, path)
	})
	if got > budget {
		t.Errorf("scanning CasbinExclude allocates %.0f times, budget is %d\n"+
			"casbin's util.KeyMatch2 costs about 2566 here; use mycasbin.KeyMatch2",
			got, budget)
	}
}

// BenchmarkCasbinExcludeScan reports what the scan adds to a request.
func BenchmarkCasbinExcludeScan(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = excludedFromCasbin("GET", "/api/v1/dept")
		}
	})
}
