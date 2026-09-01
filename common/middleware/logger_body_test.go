package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// serveWithLogger runs one request through the logger middleware and returns
// what the handler saw, with logger.enableddb set as given.
func serveWithLogger(t testing.TB, enabledDB bool, method, body string) string {
	t.Helper()

	prev := config.LoggerConfig.EnabledDB
	config.LoggerConfig.EnabledDB = enabledDB
	t.Cleanup(func() { config.LoggerConfig.EnabledDB = prev })

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerToFile())

	var seen string
	handler := func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Errorf("handler could not read the body: %v", err)
		}
		seen = string(b)
		c.Status(http.StatusOK)
	}
	r.Handle(method, "/probe", handler)

	req := httptest.NewRequest(method, "/probe", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)
	return seen
}

// The middleware rewrites Request.Body so it can log the parameters. Whatever
// else it does, the handler has to receive the request the client sent - all
// of it, whether or not the operation log is on, and whether or not the body
// is longer than what gets logged.
func TestHandlerStillSeesTheWholeBody(t *testing.T) {
	cases := []struct {
		name      string
		enabledDB bool
		body      string
	}{
		{"log off, short body", false, `{"username":"admin"}`},
		{"log on, short body", true, `{"username":"admin"}`},
		{"log off, empty body", false, ""},
		{"log on, empty body", true, ""},
		// Longer than operParamLimit: the logged copy is truncated, the body is not.
		{"log on, body past the limit", true, strings.Repeat("x", operParamLimit+4096)},
		{"log off, body past the limit", false, strings.Repeat("y", operParamLimit+4096)},
		// Exactly at the boundary, where a fencepost error would show.
		{"log on, body exactly at the limit", true, strings.Repeat("z", operParamLimit)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
				if got := serveWithLogger(t, c.enabledDB, method, c.body); got != c.body {
					t.Errorf("%s: handler saw %d bytes, the client sent %d",
						method, len(got), len(c.body))
				}
			}
		})
	}
}

// The body is read for one reason - operParam on the operation log row - and
// that row is only written when logger.enableddb is on. With it off, reading
// the body is a copy of every request made and thrown away, and a file upload
// is a POST like any other: 16MB of upload allocated about 67MB here.
//
// Allocation counts are deterministic across machines; wall-clock is not.
func TestBodyIsNotCopiedWhenTheOperationLogIsOff(t *testing.T) {
	const size = 1 << 20
	body := strings.Repeat("x", size)

	prev := config.LoggerConfig.EnabledDB
	config.LoggerConfig.EnabledDB = false
	t.Cleanup(func() { config.LoggerConfig.EnabledDB = prev })

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerToFile())
	r.POST("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })

	payload := []byte(body)
	run := func() {
		req := httptest.NewRequest(http.MethodPost, "/probe", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	var before, after uint64
	before = heapAllocs()
	run()
	after = heapAllocs()

	// The handler never reads the body, so a request that does not copy it
	// should allocate far less than the body's size. The old middleware
	// allocated about four times the body.
	if grew := after - before; grew > size/2 {
		t.Errorf("a %d-byte request allocated %d bytes with the operation log off; "+
			"the body should not be read when nothing consumes it", size, grew)
	}
}

func heapAllocs() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}
