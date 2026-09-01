package main

import (
	"fmt"
	"sort"
)

// Severity decides whether a finding stops CI.
//
// Two levels and no more. A third would immediately be used to park findings
// nobody intends to fix, and the point of this tool is that everything it
// reports is something that fails without saying so.
type Severity int

const (
	// Warn prints and does not affect the exit code.
	Warn Severity = iota
	// Error prints and makes the run fail.
	Error
)

func (s Severity) String() string {
	if s == Error {
		return "ERROR"
	}
	return "WARN"
}

// Finding is one problem, located precisely enough to open the file at it.
type Finding struct {
	Check    string   `json:"check"`
	Severity string   `json:"severity"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Col      int      `json:"col"`
	Message  string   `json:"message"`
	Related  []string `json:"related,omitempty"`

	severity Severity
}

func (f Finding) String() string {
	s := fmt.Sprintf("%s:%d:%d: [%s] %s: %s", f.File, f.Line, f.Col, f.Severity, f.Check, f.Message)
	for _, r := range f.Related {
		s += "\n    " + r
	}
	return s
}

// sortFindings orders by file, then position, then check, so two runs over the
// same tree print the same thing and a diff of the output is meaningful.
func sortFindings(fs []Finding) {
	sort.Slice(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Col != b.Col {
			return a.Col < b.Col
		}
		if a.Check != b.Check {
			return a.Check < b.Check
		}
		return a.Message < b.Message
	})
}

func hasError(fs []Finding) bool {
	for _, f := range fs {
		if f.severity == Error {
			return true
		}
	}
	return false
}
