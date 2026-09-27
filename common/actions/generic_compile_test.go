package actions_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A model paired with a request that builds a different model does not
// compile - the thing the older actions left to a type assertion at run
// time. Checked by building a package that does it, next to a control that
// is the same package with the pair corrected.
func TestGenericActionsRefuseAMismatchedPair(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	vet := func(pkg string) (string, error) {
		cmd := exec.Command("go", "vet", "./common/actions/testdata/"+pkg)
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	if out, err := vet("matched"); err != nil {
		t.Fatalf("the control does not build, so the mismatch below proves nothing:\n%s", out)
	}
	out, err := vet("mismatch")
	if err == nil {
		t.Fatal("a model paired with another model's request compiled")
	}
	if !strings.Contains(out, "does not satisfy") || !strings.Contains(out, "ToModel") {
		t.Fatalf("the build failed, but not on the pairing:\n%s", out)
	}
}
