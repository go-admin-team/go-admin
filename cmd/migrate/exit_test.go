package migrate

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// A deployment decides whether to start the new version on this command's
// exit code. Before this batch the only failure that produced one was a
// failing migration function, and it produced it by ending the process from
// inside the migration engine; moving that out would have taken the last
// reported failure with it.
func TestExitOnErrorEndsTheCommandNonZero(t *testing.T) {
	var codes []int
	osExit = func(c int) { codes = append(codes, c) }
	t.Cleanup(func() { osExit = origExit })

	var out bytes.Buffer
	exitOnError(&out, errors.New("the tenant database is unreachable"))

	if len(codes) != 1 || codes[0] != 1 {
		t.Errorf("exit codes = %v, want [1]", codes)
	}
	if !strings.Contains(out.String(), "the tenant database is unreachable") {
		t.Errorf("the reason was not reported: %q", out.String())
	}
}

func TestExitOnErrorLetsSuccessThrough(t *testing.T) {
	var codes []int
	osExit = func(c int) { codes = append(codes, c) }
	t.Cleanup(func() { osExit = origExit })

	var out bytes.Buffer
	exitOnError(&out, nil)

	if len(codes) != 0 {
		t.Errorf("a successful migration exited with %v", codes)
	}
	if out.Len() != 0 {
		t.Errorf("a successful migration wrote %q", out.String())
	}
}
