package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// heuristicNote is appended to every menu-name finding.
//
// Required by the acceptance criteria, and for a reason worth restating: this
// is the one check that compares two repositories through a regular expression,
// so it will occasionally be wrong. Reporting it as a certainty is how a tool
// gets a reputation and then an ignore comment on every line.
const heuristicNote = "\n    (heuristic: matched across repositories by regular expression - false positives are possible, verify before changing anything)"

var (
	// Vue 3 <script setup>, which is the prevailing style in go-admin-ui.
	defineOptionsName = regexp.MustCompile(`defineOptions\(\s*\{[^}]*?name:\s*['"` + "`" + `]([^'"` + "`" + `]+)['"` + "`" + `]`)
	// Options API, still present in older views.
	exportDefaultName = regexp.MustCompile(`export default\s*(?:defineComponent\(\s*)?\{[^}]*?name:\s*['"` + "`" + `]([^'"` + "`" + `]+)['"` + "`" + `]`)
)

// checkMenuNames compares the menu_name a migration seeds against the name the
// component declares.
//
// keep-alive caches by component name, and the page cache is configured by menu
// name. When they disagree nothing breaks and nothing is logged - the page just
// stops being cached, or the wrong page is evicted, and it looks like a
// performance quirk.
//
// WARN, not ERROR, and deliberately so. The two sides are in different
// repositories and different languages, so this is a regular expression against
// a .vue file, not a parse. It is scheduled to become an ERROR after two release
// cycles with no false positive; until then it must not decide an exit code.
func checkMenuNames(s *snapshot, uiDir string) ([]Finding, error) {
	info, err := os.Stat(uiDir)
	if err != nil {
		return nil, fmt.Errorf("ui dir %s: %w", uiDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("ui dir %s is not a directory", uiDir)
	}

	var out []Finding
	for _, sf := range s.Files {
		forEachStructLiteral(sf, func(lit structLiteral) {
			if !s.isMenuModel(lit) {
				return
			}
			menuType, _ := stringField(sf, lit, "MenuType")
			// A directory has no page of its own; its component is the layout.
			if menuType == "M" || menuType == "F" {
				return
			}
			component, ok := stringField(sf, lit, "Component")
			if !ok || component == "" || component == "Layout" {
				return
			}
			menuName, ok := stringField(sf, lit, "MenuName")
			if !ok || menuName == "" {
				return
			}

			path := componentFile(uiDir, component)
			if path == "" {
				return
			}
			b, err := os.ReadFile(path)
			if err != nil {
				// A component the frontend does not carry is F2's business -
				// the placeholder that says the app is not installed - not this
				// check's.
				return
			}
			declared, ok := componentName(string(b))
			if !ok {
				// No literal name at all: the build may derive one from the
				// file path. Nothing to compare, so nothing to say.
				return
			}
			if declared == menuName {
				return
			}
			expr, _ := field(lit.Lit, "MenuName")
			rel := path
			if r, err := filepath.Rel(uiDir, path); err == nil {
				rel = r
			}
			out = append(out, s.finding(Warn, checkMenuName, sf, expr,
				"menu_name %q does not match the component name %q declared in %s;\n"+
					"    keep-alive caches by component name and the cache is configured by menu name, so the page silently stops being cached.%s",
				menuName, declared, rel, heuristicNote))
		})
	}
	return out, nil
}

func stringField(sf *sourceFile, lit structLiteral, name string) (string, bool) {
	expr, ok := field(lit.Lit, name)
	if !ok {
		return "", false
	}
	return stringValue(sf, expr)
}

// componentFile maps a menu component path to a file in the frontend tree.
//
// Two roots, because F2 adds a second: anything under apps/ is an installed
// app's view, everything else is a view of the main repository.
func componentFile(uiDir, component string) string {
	c := strings.TrimPrefix(filepath.ToSlash(component), "/")
	if c == "" {
		return ""
	}
	if strings.HasPrefix(c, "apps/") {
		return filepath.Join(uiDir, filepath.FromSlash(c)+".vue")
	}
	return filepath.Join(uiDir, "views", filepath.FromSlash(c)+".vue")
}

func componentName(src string) (string, bool) {
	for _, re := range []*regexp.Regexp{defineOptionsName, exportDefaultName} {
		if m := re.FindStringSubmatch(src); m != nil {
			return m[1], true
		}
	}
	return "", false
}
