package migrate

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"go-admin/cmd/migrate/migration"
)

func at(s string) *time.Time {
	t, err := time.Parse(applyTimeLayout, s)
	if err != nil {
		panic(err)
	}
	return &t
}

// The order Status returns: version strings sorted as ASCII.
func sampleEntries() []migration.StatusEntry {
	return []migration.StatusEntry{
		{Version: "1786700001000", AppCode: "", Registered: true, Applied: true, ApplyTime: at("2026-08-20 10:00:00")},
		{Version: "1786700005000", AppCode: "", Registered: true},
		{Version: "crm-1786800001000", AppCode: "crm", Registered: true, Applied: true, ApplyTime: at("2026-08-25 14:03:11")},
		{Version: "crm-1786800002000", AppCode: "crm", Registered: true},
	}
}

func TestPrintStatusGroupsByApp(t *testing.T) {
	var buf bytes.Buffer
	if err := printStatus(&buf, sampleEntries(), ""); err != nil {
		t.Fatal(err)
	}
	got := buf.String()

	for _, want := range []string{
		"[core]",
		"[crm]",
		"applied   1786700001000      2026-08-20 10:00:00",
		"pending   1786700005000",
		"applied   crm-1786800001000  2026-08-25 14:03:11",
		"pending   crm-1786800002000",
		"2 applied, 2 pending across 2 app(s)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
	// The framework heads the list, because that is the order a full run
	// executes in.
	if strings.Index(got, "[core]") > strings.Index(got, "[crm]") {
		t.Errorf("core is not listed first:\n%s", got)
	}
}

// A row nobody registers any more is neither applied-and-current nor pending.
// Calling it applied would say the migration is in this binary, which is what
// sends someone looking for a file that was deleted.
func TestPrintStatusMarksOrphanedRows(t *testing.T) {
	entries := append(sampleEntries(), migration.StatusEntry{
		Version: "gone-1786800000000", AppCode: "gone", Applied: true, ApplyTime: at("2026-08-01 09:00:00"),
	})
	var buf bytes.Buffer
	if err := printStatus(&buf, entries, ""); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "orphaned  gone-1786800000000") {
		t.Errorf("orphaned row not marked:\n%s", got)
	}
	if !strings.Contains(got, "nothing in this binary registers them") {
		t.Errorf("orphaned rows need an explanation:\n%s", got)
	}
	if !strings.Contains(got, "2 applied, 2 pending") {
		t.Errorf("orphaned rows must not be counted as applied:\n%s", got)
	}
}

func TestPrintStatusFiltersByApp(t *testing.T) {
	var buf bytes.Buffer
	if err := printStatus(&buf, sampleEntries(), "crm"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if strings.Contains(got, "[core]") {
		t.Errorf("--app crm listed the framework:\n%s", got)
	}
	if !strings.Contains(got, "across 1 app(s)") {
		t.Errorf("output = %s", got)
	}
}

// status prints [core]; --app core has to mean the same thing.
func TestPrintStatusAppCoreSelectsTheFramework(t *testing.T) {
	var buf bytes.Buffer
	if err := printStatus(&buf, sampleEntries(), migration.FrameworkAppCode); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if strings.Contains(got, "[crm]") {
		t.Errorf("--app core listed crm:\n%s", got)
	}
	if !strings.Contains(got, "[core]") {
		t.Errorf("--app core listed nothing:\n%s", got)
	}
}

func TestPrintStatusOnAnEmptyRegistry(t *testing.T) {
	var buf bytes.Buffer
	if err := printStatus(&buf, nil, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no migrations registered and none recorded") {
		t.Errorf("output = %s", buf.String())
	}
}

func TestPrintPendingListsOnlyPendingInOrder(t *testing.T) {
	var buf bytes.Buffer
	if err := printPending(&buf, sampleEntries(), ""); err != nil {
		t.Fatal(err)
	}
	got := buf.String()

	if !strings.Contains(got, "dry-run: nothing will be written") {
		t.Errorf("dry-run must say it writes nothing:\n%s", got)
	}
	if strings.Contains(got, "1786700001000\n") || strings.Contains(got, "crm-1786800001000") {
		t.Errorf("dry-run listed already applied migrations:\n%s", got)
	}
	if !strings.Contains(got, "[core]  1786700005000") || !strings.Contains(got, "[crm]   crm-1786800002000") {
		t.Errorf("dry-run is missing pending migrations:\n%s", got)
	}
	if !strings.Contains(got, "2 migration(s) pending") {
		t.Errorf("output = %s", got)
	}
	if strings.Index(got, "1786700005000") > strings.Index(got, "crm-1786800002000") {
		t.Errorf("dry-run order does not match run order:\n%s", got)
	}
}

// An orphaned row is applied and unregistered; a dry run must not offer to
// apply it, because a real run cannot.
func TestPrintPendingSkipsOrphanedRows(t *testing.T) {
	entries := []migration.StatusEntry{
		{Version: "gone-1786800000000", AppCode: "gone", Applied: true, ApplyTime: at("2026-08-01 09:00:00")},
	}
	var buf bytes.Buffer
	if err := printPending(&buf, entries, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "nothing to apply") {
		t.Errorf("output = %s", buf.String())
	}
}

func TestPrintPendingFiltersByApp(t *testing.T) {
	var buf bytes.Buffer
	if err := printPending(&buf, sampleEntries(), "CRM"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if strings.Contains(got, "[core]") {
		t.Errorf("--app CRM listed the framework:\n%s", got)
	}
	if !strings.Contains(got, "1 migration(s) pending") {
		t.Errorf("output = %s", got)
	}
}
