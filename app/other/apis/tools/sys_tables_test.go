package tools

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// newImportRequest builds the request an import arrives in. Where the table
// list sits -- query or body -- is exactly what these tests are about, and it
// is net/http's form parsing that decides what a handler can reach, so these go
// through a real *http.Request rather than a hand-built one.
func newImportRequest(target, contentType, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return req
}

func TestTablesToImport(t *testing.T) {
	cases := []struct {
		name        string
		target      string
		contentType string
		body        string
		want        []string
		wantErr     bool
	}{{
		name:   "from the query, as every front end before v3.2.0 sent it",
		target: "/sys/tables/info?tables=sys_user,sys_post",
		want:   []string{"sys_user", "sys_post"},
	}, {
		name:        "from a JSON body, as go-admin-ui v3.2.0 sends it",
		target:      "/sys/tables/info",
		contentType: "application/json",
		body:        `{"tables":"sys_user,sys_post"}`,
		want:        []string{"sys_user", "sys_post"},
	}, {
		name:        "the query wins when a request carries both",
		target:      "/sys/tables/info?tables=sys_user",
		contentType: "application/json",
		body:        `{"tables":"sys_post"}`,
		want:        []string{"sys_user"},
	}, {
		name:   "blank entries are dropped rather than imported as a nameless table",
		target: "/sys/tables/info?tables=sys_user,,%20,sys_post",
		want:   []string{"sys_user", "sys_post"},
	}, {
		name:        "a body carrying an empty list is an error",
		target:      "/sys/tables/info",
		contentType: "application/json",
		body:        `{"tables":""}`,
		wantErr:     true,
	}, {
		name:        "a body that is not JSON at all is an error, not a panic",
		target:      "/sys/tables/info",
		contentType: "application/json",
		body:        "sys_user",
		wantErr:     true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = newImportRequest(tc.target, tc.contentType, tc.body)

			got, err := tablesToImport(c)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				if err.Error() != emptyTableNameMsg {
					t.Fatalf("message should be the one the front end shows, got %q", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// insertMsg runs one import through the wired handler. It asserts nothing about
// the import succeeding -- it cannot, over sqlite -- only about how far the
// request got, which the empty-list message is what distinguishes.
func insertMsg(t *testing.T, target, contentType, body string) bodyOf {
	t.Helper()
	return serveJSON(t,
		newEngine(t, http.MethodPost, "/sys/tables/info", SysTable{}.Insert),
		newImportRequest(target, contentType, body))
}

func TestInsert_ReadsTheTableListFromEitherPlace(t *testing.T) {
	for _, tc := range []struct {
		name        string
		target      string
		contentType string
		body        string
	}{
		{"query", "/sys/tables/info?tables=sys_user", "", ""},
		{"JSON body", "/sys/tables/info", "application/json", `{"tables":"sys_user"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := insertMsg(t, tc.target, tc.contentType, tc.body); got.Msg == emptyTableNameMsg {
				t.Fatalf("request carried a table name and was still rejected as empty: %+v", got)
			}
		})
	}
}

func TestInsert_RejectsAMissingTableList(t *testing.T) {
	if got := insertMsg(t, "/sys/tables/info", "", ""); got.Msg != emptyTableNameMsg {
		t.Fatalf("missing table list should be rejected, got %+v", got)
	}
}
