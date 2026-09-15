package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// Acceptance 14b, second half: the WARN check on its own leaves the exit code
// at zero. If it did not, the first person it misjudged would silence it, and a
// silenced check is worse than none - it looks like coverage.
func TestOnlyWarningsExitZero(t *testing.T) {
	root, ui := menuNameFixture(t, "DemoProduct", "Product")

	var buf bytes.Buffer
	code, err := run(&buf, root, options{UIDir: ui}, false)
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Errorf("exit code = %d, want 0:\n%s", code, buf.String())
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("the warning was not printed:\n%s", out)
	}
	if !strings.Contains(out, "0 error(s), 1 warning(s)") {
		t.Errorf("summary = %s", out)
	}
	if !strings.Contains(out, "Warnings do not affect the exit code.") {
		t.Errorf("output must say warnings are not fatal:\n%s", out)
	}
}

// Acceptance 14b, first half: any of the other five fails the run.
func TestAnyErrorExitsNonZero(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/models/user.go": "package models\n\ntype SysUser struct{}\n",
		"common/middleware/auth.go": `package middleware

import "go-admin/app/admin/models"

var _ = models.SysUser{}
`,
	})

	var buf bytes.Buffer
	code, err := run(&buf, root, options{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Errorf("exit code = %d, want 1:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("output = %s", buf.String())
	}
}

func TestCleanTreeExitsZero(t *testing.T) {
	root := fixture(t, map[string]string{
		"common/models/by.go": "package models\n\ntype ControlBy struct{}\n",
	})
	var buf bytes.Buffer
	code, err := run(&buf, root, options{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Errorf("exit code = %d, want 0:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "0 error(s), 0 warning(s)") {
		t.Errorf("output = %s", buf.String())
	}
}

// Every message has to name a file and a line, or the report is a puzzle rather
// than a finding.
func TestTextOutputLocatesEveryFinding(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/models/user.go": "package models\n\ntype SysUser struct{}\n",
		"common/middleware/auth.go": `package middleware

import "go-admin/app/admin/models"

var _ = models.SysUser{}
`,
	})
	var buf bytes.Buffer
	if _, err := run(&buf, root, options{}, false); err != nil {
		t.Fatal(err)
	}
	line := strings.SplitN(buf.String(), "\n", 2)[0]
	if !strings.HasPrefix(line, "common/middleware/auth.go:3:") {
		t.Errorf("first line does not locate the finding: %q", line)
	}
}

func TestJSONOutput(t *testing.T) {
	root := fixture(t, map[string]string{
		"app/admin/models/user.go": "package models\n\ntype SysUser struct{}\n",
		"common/middleware/auth.go": `package middleware

import "go-admin/app/admin/models"

var _ = models.SysUser{}
`,
	})
	var buf bytes.Buffer
	code, err := run(&buf, root, options{}, true)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Errorf("exit code = %d", code)
	}
	var findings []Finding
	if err = json.Unmarshal(buf.Bytes(), &findings); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, buf.String())
	}
	if len(findings) != 1 || findings[0].Check != checkImportBoundary || findings[0].Severity != "ERROR" {
		t.Errorf("findings = %+v", findings)
	}
}

// The tool runs on this repository in CI, so it has to be clean here. This also
// covers acceptance 13b: nothing under common/ imports app/ any more.
func TestThisRepositoryIsClean(t *testing.T) {
	var buf bytes.Buffer
	code, err := run(&buf, "../..", options{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Errorf("checksilent reports problems in this repository:\n%s", buf.String())
	}
}
