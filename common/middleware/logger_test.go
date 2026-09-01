package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/common/global"
)

// The operation-log consumer reads these keys by name off the queue message.
// Losing one costs a column in sys_opera_log and reports nothing - the request
// still succeeds, the log row is just wrong.
//
// This locks the set down across the move of the status constants out of
// app/admin/service/dto, which touched every request path.
var operaLogKeys = []string{
	"_fullPath", "operUrl", "operIp", "operLocation", "operName",
	"requestMethod", "operParam", "operTime", "jsonResult", "latencyTime",
	"statusCode", "userAgent", "createBy", "updateBy", "status",
}

func TestOperaLogFieldsAreComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sys-user", nil)
	c.Request.Header.Set("User-Agent", "go-test")

	l := operaLogFields(c, "127.0.0.1", http.StatusOK, "/api/v1/sys-user", http.MethodPost,
		12*time.Millisecond, `{"a":1}`, `{"code":200}`, http.StatusOK)

	for _, k := range operaLogKeys {
		if _, ok := l[k]; !ok {
			t.Errorf("operation log is missing %q", k)
		}
	}
	if len(l) != len(operaLogKeys) {
		t.Errorf("operation log has %d fields, expected %d; update operaLogKeys deliberately, not to make this pass",
			len(l), len(operaLogKeys))
	}
	if got := l["operUrl"]; got != "/api/v1/sys-user" {
		t.Errorf("operUrl = %v", got)
	}
	if got := l["userAgent"]; got != "go-test" {
		t.Errorf("userAgent = %v", got)
	}
}

// status is what tells a failed request from a successful one in the log table.
// It is a string, and it is the one field whose source package changed.
func TestOperaLogStatusMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name   string
		status int
		want   string
	}{
		{"ok", http.StatusOK, global.OperaStatusEnabled},
		{"error", http.StatusInternalServerError, global.OperaStatusDisabled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			l := operaLogFields(c, "127.0.0.1", tc.status, "/", http.MethodGet, 0, "", "", tc.status)
			if l["status"] != tc.want {
				t.Fatalf("status = %v, want %v", l["status"], tc.want)
			}
		})
	}
}
