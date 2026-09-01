// Command checksilent reports the failures in this repository that do not
// announce themselves: no error, no log line, behaviour quietly wrong.
//
// Six checks, five of them ERROR and one WARN. An ERROR fails the run; a WARN
// prints and does not. The split is not about how bad the consequence is - all
// six are bad - but about how certain the detection is. Everything reported as
// an ERROR is decided from this repository's own syntax. The one WARN compares
// against a second repository through a regular expression, and a check that
// can be wrong must not be able to stop a build, or the first response to it
// will be an ignore comment.
//
// Usage:
//
//	go run ./tools/checksilent
//	go run ./tools/checksilent -ui-dir ../go-admin-ui/src
//	go run ./tools/checksilent -json
//
// Or through the Makefile, which is what CI runs:
//
//	make checksilent
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var (
		root   = flag.String("root", ".", "repository root to scan")
		uiDir  = flag.String("ui-dir", "", "go-admin-ui src directory; enables the menu-name check, which is skipped without it")
		asJSON = flag.Bool("json", false, "print findings as JSON")
	)
	flag.Parse()

	code, err := run(os.Stdout, *root, options{UIDir: *uiDir}, *asJSON)
	if err != nil {
		fmt.Fprintln(os.Stderr, "checksilent:", err)
		os.Exit(2)
	}
	os.Exit(code)
}

// run returns the process exit code: 0 when nothing or only warnings were
// found, 1 when at least one error was.
func run(w io.Writer, root string, opt options, asJSON bool) (int, error) {
	s, err := load(root)
	if err != nil {
		return 0, err
	}
	findings, err := runChecks(s, opt)
	if err != nil {
		return 0, err
	}

	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err = enc.Encode(findings); err != nil {
			return 0, err
		}
	} else {
		for _, f := range findings {
			fmt.Fprintln(w, f)
		}
		printSummary(w, findings, opt, s)
	}

	if hasError(findings) {
		return 1, nil
	}
	return 0, nil
}

func printSummary(w io.Writer, findings []Finding, opt options, s *snapshot) {
	var errors, warnings int
	for _, f := range findings {
		if f.severity == Error {
			errors++
		} else {
			warnings++
		}
	}
	if len(findings) > 0 {
		fmt.Fprintln(w)
	}
	fmt.Fprintf(w, "checksilent: %d error(s), %d warning(s)\n", errors, warnings)
	if warnings > 0 {
		fmt.Fprintln(w, "Warnings do not affect the exit code.")
	}
	if opt.UIDir == "" {
		fmt.Fprintf(w, "The %s check was skipped: pass -ui-dir <go-admin-ui>/src to run it.\n", checkMenuName)
	}
	// Said out loud so nobody reads a clean run as "the boundary holds
	// everywhere it was declared". core/ is a separate module and has no
	// directory here, so this repository's copy of the check cannot cover it.
	if scanned, absent := ScannedContractRoots(s); len(absent) > 0 {
		fmt.Fprintf(w, "The %s check covered %s; %s does not exist here and was not scanned.\n",
			checkImportBoundary, strings.Join(scanned, ", "), strings.Join(absent, ", "))
	}
}
