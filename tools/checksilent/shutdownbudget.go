package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

// The files this check compares, and the package the fallbacks come from.
const (
	settingsFile     = "config/settings.yml"
	k8sDeployFile    = "scripts/k8s/deploy.yml"
	pkgHostConfig    = "config"
	drainConstName   = "DefaultDrainSeconds"
	serverConstName  = "DefaultServerSeconds"
	cleanupConstName = "DefaultCleanupSeconds"
)

// graceMarginSeconds is the headroom a shutdown budget needs beyond itself.
//
// Spelled once and used by both checks below, because they fail the same way:
// somebody raises a budget in config/settings.yml and does not go looking for
// the two other places that have to allow room for it. Two margins would
// eventually be two different numbers.
const graceMarginSeconds = 5

// checkShutdownBudgetFitsGrace compares the shutdown budget this repository
// ships against the stop grace period its own Kubernetes manifest allows.
//
// The two are not merely adjacent examples. scripts/k8s/prerun.sh builds the
// settings-admin ConfigMap out of config/settings.yml, and the Deployment
// mounts that ConfigMap - so the manifest deploys that file.
//
// The budgets are spent one after the other, and when their sum reaches
// terminationGracePeriodSeconds the kubelet sends SIGKILL while the cleanup
// callbacks are still running. Nothing reports it: the pod disappears
// mid-shutdown and it reads as a crash rather than as a number that was raised
// in one file and not the other. Which is how it would be raised - drain is
// the interesting knob and the grace period is in a different directory.
//
// Two levels, and an overrun is not also reported as a shortage of headroom:
// every Error satisfies the Warn condition too, and an Error that always drags
// a duplicate Warn behind it teaches people to skip Warns.
//
// A preStop hook counts, even though the shipped manifest has none. It is
// spent before the process is told anything, so it is added to the budget
// rather than overlapping it - and a self-check that cannot see it would
// understate the real cost by however long somebody set it to, which is worse
// than not checking.
//
// It reports nothing when either file is absent and when the manifest sets no
// grace period, because there is then no second number to disagree with.
func checkShutdownBudgetFitsGrace(s *snapshot) ([]Finding, error) {
	budget, ok, err := shippedShutdownBudget(s)
	if err != nil || !ok {
		return nil, err
	}
	m, ok, err := readManifest(s)
	if err != nil || !ok {
		return nil, err
	}
	if m.grace == nil {
		return nil, nil
	}

	var out []Finding
	if m.preStopUnreadable {
		out = append(out, Finding{
			Check:    checkShutdownGrace,
			Severity: Warn.String(),
			File:     k8sDeployFile,
			Line:     m.preStopLine,
			Col:      1,
			Message: "this preStop hook is not a sleep, so how long it takes cannot be read here " +
				"and is not in the sum below; it is spent before the process is told anything, " +
				"so whatever it costs has to fit inside terminationGracePeriodSeconds as well.",
			severity: Warn,
		})
	}

	total := m.preStop + budget.drain + budget.server + budget.cleanup
	grace := *m.grace
	spelled := fmt.Sprintf("preStop %d + drain %d + server %d + cleanup %d",
		m.preStop, budget.drain, budget.server, budget.cleanup)

	switch {
	case total >= grace:
		out = append(out, Finding{
			Check:    checkShutdownGrace,
			Severity: Error.String(),
			File:     k8sDeployFile,
			Line:     m.graceLine,
			Col:      1,
			Message: fmt.Sprintf(
				"terminationGracePeriodSeconds is %d and the shutdown takes %d (%s, from %s); "+
					"SIGKILL would arrive while the cleanup callbacks are still running. "+
					"Raise it to %d, or take %d off the budget.",
				grace, total, spelled, settingsFile,
				total+graceMarginSeconds, total+graceMarginSeconds-grace),
			severity: Error,
		})
	case total+graceMarginSeconds > grace:
		out = append(out, Finding{
			Check:    checkShutdownGrace,
			Severity: Warn.String(),
			File:     k8sDeployFile,
			Line:     m.graceLine,
			Col:      1,
			Message: fmt.Sprintf(
				"terminationGracePeriodSeconds is %d and the shutdown takes %d (%s, from %s), "+
					"which leaves under %ds of headroom; a callback that runs slightly long is "+
					"cut off. Raise it to %d.",
				grace, total, spelled, settingsFile, graceMarginSeconds, total+graceMarginSeconds),
			severity: Warn,
		})
	}
	return out, nil
}

// dockerStopArgs matches a stop command in a script or a workflow.
var (
	dockerStopArgs = regexp.MustCompile(`\bdocker\s+stop\b`)
	// --timeout is the current name, --time its deprecated spelling and -t the
	// short form; docker still accepts all three, so all three are read. The
	// long name comes first because --time is a prefix of it, and a flag that
	// the check cannot read is reported as no deadline at all - which would
	// have this tool pressing people towards the deprecated spelling.
	dockerStopTime = regexp.MustCompile(`(--timeout|--time|-t)[=\s]*(\d+)`)
)

// checkDockerStopGrace reports a stop path that does not allow this process
// the time it spends shutting down.
//
// docker allows ten seconds unless told otherwise, and that number is nowhere
// near the command - so a budget raised in config/settings.yml passes every
// test, deploys, and then has its cleanup callbacks killed on the next
// release. Same failure as the manifest's grace period, same arithmetic, same
// margin; only the file it lives in is different.
//
// Both ways of stopping this repository's container are covered, because
// covering one of two identical paths is what produces a clean run that means
// nothing: `docker stop` in a workflow or a script, and stop_grace_period in
// the compose file the Makefile's own `make run` uses.
//
// An absent deadline is reported rather than assumed to be ten: the value that
// applies is then invisible at the call site and cannot follow the budget.
func checkDockerStopGrace(s *snapshot) ([]Finding, error) {
	budget, ok, err := shippedShutdownBudget(s)
	if err != nil || !ok {
		return nil, err
	}
	total := budget.drain + budget.server + budget.cleanup
	spelled := fmt.Sprintf("drain %d + server %d + cleanup %d",
		budget.drain, budget.server, budget.cleanup)

	sites, err := findStopDeadlines(s)
	if err != nil {
		return nil, err
	}

	var out []Finding
	for _, site := range sites {
		finding := Finding{
			Check: checkDockerStop,
			File:  site.file,
			Line:  site.line,
			Col:   1,
		}
		switch {
		case !site.set:
			finding.Severity, finding.severity = Error.String(), Error
			finding.Message = fmt.Sprintf(
				"%s allows the default %d seconds, and this process spends %d shutting down "+
					"(%s, from %s). %s.",
				site.what, dockerDefaultGraceSeconds, total, spelled, settingsFile,
				site.fix(total+graceMarginSeconds))
		case site.seconds <= total:
			finding.Severity, finding.severity = Error.String(), Error
			finding.Message = fmt.Sprintf(
				"%s allows %d seconds and this shutdown takes %d (%s, from %s); the cleanup "+
					"callbacks are killed part-way through. %s.",
				site.what, site.seconds, total, spelled, settingsFile,
				site.fix(total+graceMarginSeconds))
		case site.seconds < total+graceMarginSeconds:
			finding.Severity, finding.severity = Warn.String(), Warn
			finding.Message = fmt.Sprintf(
				"%s allows %d seconds over a shutdown that takes %d (%s, from %s), which leaves "+
					"under %ds of headroom. %s.",
				site.what, site.seconds, total, spelled, settingsFile, graceMarginSeconds,
				site.fix(total+graceMarginSeconds))
		default:
			continue
		}
		out = append(out, finding)
	}
	return out, nil
}

// dockerDefaultGraceSeconds is what docker allows a container to stop in when
// nothing says otherwise. It applies to `docker stop` and to compose alike.
const dockerDefaultGraceSeconds = 10

// stopSite is one place this repository decides how long a container gets.
type stopSite struct {
	file string
	line int
	// what names the setting in the finding, in the spelling of the file it
	// was found in.
	what string
	// compose says which of the two fixes to suggest.
	compose bool
	seconds int
	set     bool
}

func (s stopSite) fix(seconds int) string {
	if s.compose {
		return fmt.Sprintf("Set stop_grace_period: %ds", seconds)
	}
	return fmt.Sprintf("Pass --timeout %d", seconds)
}

func findStopDeadlines(s *snapshot) ([]stopSite, error) {
	sites, err := findDockerStops(s.Root)
	if err != nil {
		return nil, err
	}
	compose, err := findComposeServices(s)
	if err != nil {
		return nil, err
	}
	return append(sites, compose...), nil
}

// dockerStopExtensions and dockerStopNames are where a stop command can be
// written in this repository: workflows, shell scripts and the Makefile.
var (
	dockerStopExtensions = map[string]bool{".yml": true, ".yaml": true, ".sh": true, ".bash": true}
	dockerStopNames      = map[string]bool{"Makefile": true, "makefile": true}
)

func findDockerStops(root string) ([]stopSite, error) {
	var out []stopSite
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != root && skippedDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !dockerStopExtensions[filepath.Ext(path)] && !dockerStopNames[info.Name()] {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(b), "\n") {
			// A commented-out command is not one that runs, and the settings
			// file describes `docker stop` in prose right beside the budget
			// this check reads.
			if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, "#") {
				continue
			}
			if !dockerStopArgs.MatchString(line) {
				continue
			}
			site := stopSite{
				file: filepath.ToSlash(rel),
				line: i + 1,
				what: "`docker stop` with no --timeout",
			}
			if m := dockerStopTime.FindStringSubmatch(line); m != nil {
				seconds, err := strconv.Atoi(m[2])
				if err != nil {
					continue
				}
				// Quoted back in the spelling it was written in, so the
				// message cannot misreport what the line says.
				site.what = fmt.Sprintf("`docker stop %s %d`", m[1], seconds)
				site.seconds, site.set = seconds, true
			}
			out = append(out, site)
		}
		return nil
	})
	return out, err
}

type shutdownSeconds struct{ drain, server, cleanup int }

// shippedShutdownBudget reads extend.shutdown out of the settings file this
// repository ships, filling in whatever it leaves out from the Go constants
// that do the same at run time.
//
// Taking the fallbacks from the snapshot rather than repeating 0/5/3 here is
// what keeps this honest when the defaults move: a tool that carries its own
// copy of the number it is checking eventually checks the wrong one.
//
// A negative value is left alone. config.Shutdown.Budget refuses it and the
// server does not start, so it is not a failure that passes unnoticed - and
// adding a negative into the sums above would understate them.
func shippedShutdownBudget(s *snapshot) (shutdownSeconds, bool, error) {
	raw, ok, err := readRepoFile(s, settingsFile)
	if err != nil || !ok {
		return shutdownSeconds{}, false, err
	}

	var doc struct {
		Settings struct {
			Extend struct {
				Shutdown *struct {
					Drain   *int `yaml:"drain"`
					Server  *int `yaml:"server"`
					Cleanup *int `yaml:"cleanup"`
				} `yaml:"shutdown"`
			} `yaml:"extend"`
		} `yaml:"settings"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return shutdownSeconds{}, false, fmt.Errorf("%s: %w", settingsFile, err)
	}
	section := doc.Settings.Extend.Shutdown
	if section == nil {
		return shutdownSeconds{}, false, nil
	}

	defaults, ok := s.hostConfigDefaults()
	if !ok {
		// The constants moved or were renamed. Reporting nothing would let the
		// check go quiet, which is the failure it exists to catch, so this
		// stops the run instead.
		return shutdownSeconds{}, false, fmt.Errorf(
			"%s has extend.shutdown but package %s declares no %s/%s/%s to fall back on",
			settingsFile, pkgHostConfig, drainConstName, serverConstName, cleanupConstName)
	}

	budget := shutdownSeconds{
		drain:   orDefault(section.Drain, defaults.drain),
		server:  orDefault(section.Server, defaults.server),
		cleanup: orDefault(section.Cleanup, defaults.cleanup),
	}
	if budget.drain < 0 || budget.server < 0 || budget.cleanup < 0 {
		return shutdownSeconds{}, false, nil
	}
	return budget, true, nil
}

func orDefault(configured *int, fallback int) int {
	if configured != nil {
		return *configured
	}
	return fallback
}

// hostConfigDefaults reads the three fallback constants out of the parsed tree.
func (s *snapshot) hostConfigDefaults() (shutdownSeconds, bool) {
	for _, sf := range s.Files {
		if sf.Pkg != s.pkg(pkgHostConfig) {
			continue
		}
		drain, okDrain := sf.consts[drainConstName]
		server, okServer := sf.consts[serverConstName]
		cleanup, okCleanup := sf.consts[cleanupConstName]
		if okDrain && okServer && okCleanup {
			return shutdownSeconds{int(drain), int(server), int(cleanup)}, true
		}
	}
	return shutdownSeconds{}, false
}

// manifest is what the shipped Deployment says about how long it will wait.
type manifest struct {
	grace     *int
	graceLine int
	// preStop is the longest sleep any container's hook performs, since the
	// hooks of several containers run at the same time.
	preStop     int
	preStopLine int

	preStopUnreadable bool
}

var preStopSleep = regexp.MustCompile(`\bsleep\s+(\d+)s?\b`)

// readManifest finds the grace period and the preStop hooks in the shipped
// manifest, with the lines they are on so a finding can be opened at them.
//
// The file holds several documents and only the Deployment carries a pod
// template, so every document is decoded and the first one with a grace period
// wins.
func readManifest(s *snapshot) (manifest, bool, error) {
	raw, ok, err := readRepoFile(s, k8sDeployFile)
	if err != nil || !ok {
		return manifest{}, false, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(raw))
	for {
		var doc struct {
			Spec struct {
				Template struct {
					Spec struct {
						Grace      *int `yaml:"terminationGracePeriodSeconds"`
						Containers []struct {
							Lifecycle struct {
								// A value, not a pointer: yaml.v3 only hands
								// the raw node to a field of type yaml.Node,
								// and a *yaml.Node field is allocated and left
								// empty - which reads as "the hook is there but
								// unreadable" for every manifest that has one.
								PreStop yaml.Node `yaml:"preStop"`
							} `yaml:"lifecycle"`
						} `yaml:"containers"`
					} `yaml:"spec"`
				} `yaml:"template"`
			} `yaml:"spec"`
		}
		switch err := dec.Decode(&doc); {
		case errors.Is(err, io.EOF):
			return manifest{}, false, nil
		case err != nil:
			return manifest{}, false, fmt.Errorf("%s: %w", k8sDeployFile, err)
		}

		pod := doc.Spec.Template.Spec
		if pod.Grace == nil && len(pod.Containers) == 0 {
			continue
		}

		m := manifest{
			grace:     pod.Grace,
			graceLine: lineOf(raw, "terminationGracePeriodSeconds:"),
		}
		for _, c := range pod.Containers {
			hook := c.Lifecycle.PreStop
			if hook.Kind == 0 {
				continue
			}
			m.preStopLine = hook.Line
			if seconds, ok := preStopSeconds(&hook); ok {
				// The longest one, not the sum: the hooks of several
				// containers run at the same time.
				if seconds > m.preStop {
					m.preStop = seconds
				}
				continue
			}
			m.preStopUnreadable = true
		}
		if m.grace == nil {
			continue
		}
		return m, true, nil
	}
}

// preStopSeconds reads how long a hook sleeps for.
//
// Every scalar under the hook is joined and searched, because the sleep can be
// written as one argument or as several: ["sh","-c","sleep 10"] and
// ["sleep","10"] both wait ten seconds.
func preStopSeconds(node *yaml.Node) (int, bool) {
	var words []string
	var walk func(*yaml.Node)
	walk = func(n *yaml.Node) {
		if n == nil {
			return
		}
		if n.Kind == yaml.ScalarNode {
			words = append(words, n.Value)
		}
		for _, child := range n.Content {
			walk(child)
		}
	}
	walk(node)

	m := preStopSleep.FindStringSubmatch(strings.Join(words, " "))
	if m == nil {
		return 0, false
	}
	seconds, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, false
	}
	return seconds, true
}

// lineOf locates a key for a finding's position. A miss reports line 1 rather
// than failing: the position is where to look, and the message is the finding.
func lineOf(content []byte, key string) int {
	for i, l := range strings.Split(string(content), "\n") {
		if strings.Contains(l, key) && !strings.HasPrefix(strings.TrimSpace(l), "#") {
			return i + 1
		}
	}
	return 1
}

// readRepoFile reads a file relative to the scanned root, reporting absence
// rather than failing on it: the checks run over fixtures that carry only what
// the check under test needs.
func readRepoFile(s *snapshot, rel string) ([]byte, bool, error) {
	b, err := os.ReadFile(filepath.Join(s.Root, filepath.FromSlash(rel)))
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, false, nil
	case err != nil:
		return nil, false, err
	}
	return b, true, nil
}

// composeFiles are the names Docker Compose looks for, in its own order of
// preference.
var composeFiles = []string{"compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml"}

// composeDuration matches the durations compose accepts for
// stop_grace_period: a bare number of seconds, or hours, minutes and seconds
// in that order.
var composeDuration = regexp.MustCompile(`^(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s?)?$`)

// findComposeServices reports the stop_grace_period of every compose service
// that runs this repository's own image.
//
// Only those services. The grace period of a database or a cache alongside it
// is not this process's shutdown budget, and reporting one against the other
// would be arithmetic about two unrelated things.
func findComposeServices(s *snapshot) ([]stopSite, error) {
	var out []stopSite
	for _, name := range composeFiles {
		raw, ok, err := readRepoFile(s, name)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		var root yaml.Node
		if err := yaml.Unmarshal(raw, &root); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if len(root.Content) == 0 {
			continue
		}
		services := mapValue(root.Content[0], "services")
		if services == nil {
			continue
		}

		for i := 0; i+1 < len(services.Content); i += 2 {
			key, service := services.Content[i], services.Content[i+1]
			if !runsThisRepo(service, s.ModulePath) {
				continue
			}
			site := stopSite{
				file:    name,
				line:    key.Line,
				what:    fmt.Sprintf("service %s, which sets no stop_grace_period,", key.Value),
				compose: true,
			}
			if grace := mapValue(service, "stop_grace_period"); grace != nil {
				seconds, ok := composeSeconds(grace.Value)
				if !ok {
					// A duration this cannot read is left alone rather than
					// guessed at: compose knows what it means, and inventing a
					// number here would report against a value nobody wrote.
					continue
				}
				site.line = grace.Line
				site.what = fmt.Sprintf("stop_grace_period on service %s", key.Value)
				site.seconds, site.set = seconds, true
			}
			out = append(out, site)
		}
	}
	return out, nil
}

// runsThisRepo reports whether a compose service starts the image this
// repository builds - by building it, or by naming it.
func runsThisRepo(service *yaml.Node, modulePath string) bool {
	if mapValue(service, "build") != nil {
		return true
	}
	image := mapValue(service, "image")
	if image == nil {
		return false
	}
	repository := image.Value
	if i := strings.LastIndex(repository, ":"); i > strings.LastIndex(repository, "/") {
		repository = repository[:i]
	}
	return baseName(repository) == baseName(modulePath)
}

func baseName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}

func composeSeconds(value string) (int, bool) {
	m := composeDuration.FindStringSubmatch(strings.TrimSpace(value))
	if m == nil || m[1]+m[2]+m[3] == "" {
		return 0, false
	}
	var total int
	for i, unit := range []int{3600, 60, 1} {
		if m[i+1] == "" {
			continue
		}
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return 0, false
		}
		total += n * unit
	}
	return total, true
}

// mapValue returns the value a mapping node holds for key.
func mapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}
