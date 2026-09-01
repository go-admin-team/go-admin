package migrate

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"go-admin/cmd/migrate/migration"
)

const applyTimeLayout = "2006-01-02 15:04:05"

// printStatus lists every migration this binary knows about together with every
// row already in sys_migration, grouped by app.
//
// filter is an app code as typed on the command line; empty means every app.
func printStatus(w io.Writer, entries []migration.StatusEntry, filter string) error {
	entries = filterByApp(entries, filter)

	groups, order := groupByApp(entries)
	if len(order) == 0 {
		_, err := fmt.Fprintln(w, "no migrations registered and none recorded")
		return err
	}

	// One width for the whole listing rather than one per group: the versions
	// of two apps line up, so a long list can be read down the column.
	width := versionWidth(entries)

	var applied, pending, orphaned int
	for i, app := range order {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "[%s]\n", app)
		for _, e := range groups[app] {
			state := "pending"
			switch {
			case e.Applied && !e.Registered:
				state = "orphaned"
				orphaned++
			case e.Applied:
				state = "applied"
				applied++
			default:
				pending++
			}
			fmt.Fprintln(w, strings.TrimRight(
				fmt.Sprintf("  %-*s%-*s%s", stateWidth, state, width, e.Version, formatApplyTime(e.ApplyTime)), " "))
		}
	}

	fmt.Fprintf(w, "\n%d applied, %d pending across %d app(s)\n", applied, pending, len(order))
	if orphaned > 0 {
		fmt.Fprintf(w, "%d orphaned: recorded in sys_migration, but nothing in this binary registers them.\n"+
			"Expected after a migration file is removed or an app is uninstalled; they will not run again.\n", orphaned)
	}
	return nil
}

// printPending is --dry-run: the same data as status, narrowed to what an
// actual run would do and printed in the order it would do it.
//
// It reads and prints. Every write path - AutoMigrate on sys_migration
// included - is on the other branch in initDB, so a dry run leaves the database
// byte for byte as it found it.
func printPending(w io.Writer, entries []migration.StatusEntry, filter string) error {
	entries = filterByApp(entries, filter)

	fmt.Fprintln(w, "dry-run: nothing will be written")

	pending := make([]migration.StatusEntry, 0, len(entries))
	for _, e := range entries {
		// An orphaned row is recorded and unregistered; a real run cannot
		// apply it, so a dry run must not offer to.
		if !e.Applied && e.Registered {
			pending = append(pending, e)
		}
	}
	if len(pending) == 0 {
		_, err := fmt.Fprintln(w, "nothing to apply")
		return err
	}

	appWidth := 0
	for _, e := range pending {
		if n := len(migration.DisplayAppCode(e.AppCode)) + 2; n > appWidth {
			appWidth = n
		}
	}

	fmt.Fprintln(w, "would apply, in this order:")
	for _, e := range pending {
		fmt.Fprintf(w, "  %-*s%s\n", appWidth+2, "["+migration.DisplayAppCode(e.AppCode)+"]", e.Version)
	}
	fmt.Fprintf(w, "\n%d migration(s) pending\n", len(pending))
	return nil
}

// stateWidth is the width of the applied/pending/orphaned column, sized to the
// longest of the three plus a gap.
const stateWidth = len("orphaned") + 2

func versionWidth(entries []migration.StatusEntry) int {
	width := 0
	for _, e := range entries {
		if n := len(e.Version) + 2; n > width {
			width = n
		}
	}
	return width
}

// filterByApp keeps the entries of one app. The filter is matched after the
// same normalisation ForApp applies, so --app CRM finds crm.
func filterByApp(entries []migration.StatusEntry, filter string) []migration.StatusEntry {
	if filter == "" {
		return entries
	}
	want := migration.AppFilter(filter)
	out := make([]migration.StatusEntry, 0, len(entries))
	for _, e := range entries {
		if e.AppCode == want {
			out = append(out, e)
		}
	}
	return out
}

// groupByApp buckets entries by display name and returns the buckets plus the
// order to print them in: the framework first, then apps alphabetically. That
// is also the order a full run executes them in, because version strings sort
// as ASCII and the framework's are bare digits.
func groupByApp(entries []migration.StatusEntry) (map[string][]migration.StatusEntry, []string) {
	groups := make(map[string][]migration.StatusEntry)
	for _, e := range entries {
		app := migration.DisplayAppCode(e.AppCode)
		groups[app] = append(groups[app], e)
	}
	order := make([]string, 0, len(groups))
	for app := range groups {
		order = append(order, app)
	}
	sort.Slice(order, func(i, j int) bool {
		if (order[i] == migration.FrameworkAppCode) != (order[j] == migration.FrameworkAppCode) {
			return order[i] == migration.FrameworkAppCode
		}
		return order[i] < order[j]
	})
	return groups, order
}

func formatApplyTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(applyTimeLayout)
}
