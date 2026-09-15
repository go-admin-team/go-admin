package main

import (
	"strings"
	"testing"
)

// hostConfigSource is the part of config/extend.go this check reads: the
// fallbacks it applies to whatever the settings file leaves out.
const hostConfigSource = `package config

const (
	DefaultDrainSeconds   = 0
	DefaultServerSeconds  = 5
	DefaultCleanupSeconds = 3
)
`

// factorySettings is what this repository ships: 0 + 5 + 3.
const factorySettings = "settings:\n  extend:\n    shutdown:\n      drain: 0\n      server: 5\n      cleanup: 3\n"

func settingsWith(shutdown string) string {
	return "settings:\n  extend:\n" + shutdown
}

// deployWith builds a manifest with the given container extras and pod-level
// lines, in the shape the shipped one has.
func deployWith(containerExtra, podExtra string) string {
	return `---
apiVersion: v1
kind: Service
metadata:
  name: go-admin
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-admin-v1
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: go-admin
        image: go-admin
` + containerExtra + podExtra
}

func graceOf(seconds string) string {
	return "      terminationGracePeriodSeconds: " + seconds + "\n"
}

const preStopSleep25 = `        lifecycle:
          preStop:
            exec:
              command: ["sh", "-c", "sleep 25"]
`

// The six scenarios worked through in the technical plan, plus the one that
// only fails when preStop is left out of the sum.
//
// The values matter. "raise server to 25" gives 28, which is under a grace
// period of 30 and reaches only the WARN level - it would not show that the
// ERROR level works at all.
func TestShutdownBudgetAgainstTheGracePeriod(t *testing.T) {
	for _, tc := range []struct {
		name     string
		settings string
		deploy   string
		want     Severity
		contains string
	}{
		{
			name:     "the shipped defaults, with headroom",
			settings: factorySettings,
			deploy:   deployWith("", graceOf("30")),
			want:     -1,
		},
		{
			// The one this check exists for: drain is the interesting knob and
			// the grace period is in another directory, so raising one and not
			// the other is the natural mistake.
			name:     "the budget was raised and the manifest was not",
			settings: settingsWith("    shutdown:\n      drain: 0\n      server: 30\n      cleanup: 3\n"),
			deploy:   deployWith("", graceOf("30")),
			want:     Error,
			contains: "preStop 0 + drain 0 + server 30 + cleanup 3",
		},
		{
			// Equal is not a fit: the grace period is when SIGKILL is sent, so
			// a budget that ends exactly then leaves nothing time to return.
			name:     "the grace period was lowered to the budget",
			settings: factorySettings,
			deploy:   deployWith("", graceOf("8")),
			want:     Error,
		},
		{
			name:     "fits, but with nothing to spare",
			settings: factorySettings,
			deploy:   deployWith("", graceOf("12")),
			want:     Warn,
			contains: "leaves under 5s of headroom",
		},
		{
			// The hook is spent before the process is told anything, so it is
			// added to the budget rather than overlapping it.
			name:     "a preStop hook is part of the budget",
			settings: factorySettings,
			deploy:   deployWith(preStopSleep25, graceOf("30")),
			want:     Error,
			contains: "preStop 25 + drain 0 + server 5 + cleanup 3",
		},
		{
			// The example from the review: 10 + 10 + 5 + 3 against 30.
			name:     "preStop and a drain window together, just fitting",
			settings: settingsWith("    shutdown:\n      drain: 10\n      server: 5\n      cleanup: 3\n"),
			deploy: deployWith(`        lifecycle:
          preStop:
            exec:
              command: ["sleep", "10"]
`, graceOf("30")),
			want: Warn,
		},
		{
			name:     "no shutdown section",
			settings: settingsWith("    rateLimit:\n      inboundQPS: 200\n"),
			deploy:   deployWith("", graceOf("30")),
			want:     -1,
		},
		{
			// Nothing to disagree with. A manifest without a grace period gets
			// the Kubernetes default, which this file cannot see, and guessing
			// at it would make the check wrong rather than quiet.
			name:     "the manifest sets no grace period",
			settings: settingsWith("    shutdown:\n      drain: 300\n"),
			deploy:   deployWith("", ""),
			want:     -1,
		},
		{
			// config.Shutdown.Budget refuses this and the server does not
			// start, so it is not a failure that passes unnoticed - and adding
			// a negative into the sum would understate it.
			name:     "a negative budget is left to the run-time refusal",
			settings: settingsWith("    shutdown:\n      drain: -100\n      server: 5\n      cleanup: 3\n"),
			deploy:   deployWith("", graceOf("5")),
			want:     -1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{
				"config/extend.go":       hostConfigSource,
				"config/settings.yml":    tc.settings,
				"scripts/k8s/deploy.yml": tc.deploy,
			})
			got := only(t, check(t, root, options{}), checkShutdownGrace)
			if tc.want < 0 {
				if len(got) != 0 {
					t.Fatalf("reported %d findings, want none:\n%v", len(got), got)
				}
				return
			}
			if len(got) != 1 {
				// An ERROR also satisfies the WARN condition, so a second
				// finding here means the two levels were not made exclusive -
				// and an ERROR that always drags a duplicate WARN behind it
				// teaches people to skip WARNs.
				t.Fatalf("reported %d findings, want exactly 1:\n%v", len(got), got)
			}
			if got[0].severity != tc.want {
				t.Errorf("reported %s, want %s: %s", got[0].Severity, tc.want, got[0].Message)
			}
			if tc.contains != "" && !strings.Contains(got[0].Message, tc.contains) {
				t.Errorf("message %q does not contain %q", got[0].Message, tc.contains)
			}
			if got[0].File != k8sDeployFile {
				t.Errorf("reported against %s, want %s", got[0].File, k8sDeployFile)
			}
			if want := lineOf([]byte(tc.deploy), "terminationGracePeriodSeconds:"); got[0].Line != want {
				t.Errorf("reported line %d, want %d", got[0].Line, want)
			}
		})
	}
}

// A hook whose duration cannot be read is said out loud rather than counted as
// nothing. It is still spent inside the grace period, and a self-check that
// silently valued it at zero would be the understatement this check exists to
// prevent.
func TestAnUnreadablePreStopIsReported(t *testing.T) {
	root := fixture(t, map[string]string{
		"config/extend.go":    hostConfigSource,
		"config/settings.yml": factorySettings,
		"scripts/k8s/deploy.yml": deployWith(`        lifecycle:
          preStop:
            httpGet:
              path: /drain
              port: 8000
`, graceOf("30")),
	})
	got := only(t, check(t, root, options{}), checkShutdownGrace)
	if len(got) != 1 {
		t.Fatalf("reported %d findings, want 1:\n%v", len(got), got)
	}
	if got[0].severity != Warn {
		t.Errorf("reported %s, want WARN", got[0].Severity)
	}
	if !strings.Contains(got[0].Message, "not a sleep") {
		t.Errorf("message %q does not say why the hook could not be read", got[0].Message)
	}
}

// Either file missing means there is nothing to compare, which is the state
// every other check's fixture is in.
func TestShutdownBudgetIsSkippedWithoutBothFiles(t *testing.T) {
	for _, files := range []map[string]string{
		{"config/extend.go": hostConfigSource},
		{"config/extend.go": hostConfigSource, "config/settings.yml": settingsWith("    shutdown:\n      drain: 300\n")},
		{"config/extend.go": hostConfigSource, "scripts/k8s/deploy.yml": deployWith("", graceOf("30"))},
	} {
		root := fixture(t, files)
		if got := only(t, check(t, root, options{}), checkShutdownGrace); len(got) != 0 {
			t.Errorf("reported %d findings with only %d file(s):\n%v", len(got), len(files), got)
		}
	}
}

// A tool that cannot find the defaults it is meant to apply has to say so.
// Reporting nothing would be the failure this whole tool is about: a check
// that stops checking and goes on printing a clean run.
func TestShutdownBudgetStopsWhenTheFallbacksAreGone(t *testing.T) {
	root := fixture(t, map[string]string{
		"config/extend.go":       "package config\n\nconst DefaultDrainSeconds = 0\n",
		"config/settings.yml":    settingsWith("    shutdown:\n      drain: 1\n"),
		"scripts/k8s/deploy.yml": deployWith("", graceOf("30")),
	})
	s, err := load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, err := runChecks(s, options{}); err == nil {
		t.Fatal("runChecks succeeded with the fallback constants renamed away")
	} else if !strings.Contains(err.Error(), serverConstName) {
		t.Errorf("error %q does not name the missing constant", err)
	}
}

// The same arithmetic and the same margin as the manifest check, against the
// other place a shutdown gets cut short.
func TestDockerStopAgainstTheShutdownBudget(t *testing.T) {
	for _, tc := range []struct {
		name     string
		script   string
		want     Severity
		contains string
	}{
		{
			name:   "explicit and generous",
			script: "sudo docker stop --timeout 30 \"$PREV\"\n",
			want:   -1,
		},
		{
			// --time is the deprecated spelling of the same flag and docker
			// still honours it. A check that could not read it would report a
			// deadline that exists as missing, and push whoever fixed that
			// towards a flag that is on its way out.
			name:   "the deprecated spelling still counts",
			script: "docker stop --time 30 go-admin\n",
			want:   -1,
		},
		{
			name:   "the short form counts too",
			script: "docker stop -t 30 go-admin\n",
			want:   -1,
		},
		{
			// docker's default is 10 and this process spends 8, so it happens
			// to work today - and would stop working the first time anybody
			// configures a drain window, without the command changing.
			name:     "no deadline at all",
			script:   "sudo docker stop \"$PREV\" >/dev/null\n",
			want:     Error,
			contains: "Pass --timeout 13",
		},
		{
			name:     "shorter than the shutdown",
			script:   "docker stop --timeout 5 go-admin\n",
			want:     Error,
			contains: "allows 5 seconds and this shutdown takes 8",
		},
		{
			// Quoted back in the spelling that was written, so the message
			// cannot misreport the line it is pointing at.
			name:     "the message quotes the flag that was used",
			script:   "docker stop -t 5 go-admin\n",
			want:     Error,
			contains: "`docker stop -t 5`",
		},
		{
			name:   "longer than the shutdown but inside the margin",
			script: "docker stop --timeout=10 go-admin\n",
			want:   Warn,
		},
		{
			name:   "exactly the margin",
			script: "docker stop --timeout 13 go-admin\n",
			want:   -1,
		},
		{
			// The settings file describes `docker stop` in prose right beside
			// the budget this check reads.
			name:   "a commented-out command is not one that runs",
			script: "# docker stop go-admin\n",
			want:   -1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{
				"config/extend.go":    hostConfigSource,
				"config/settings.yml": factorySettings,
				"scripts/deploy.sh":   "#!/bin/sh\n" + tc.script,
			})
			got := only(t, check(t, root, options{}), checkDockerStop)
			if tc.want < 0 {
				if len(got) != 0 {
					t.Fatalf("reported %d findings, want none:\n%v", len(got), got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("reported %d findings, want 1:\n%v", len(got), got)
			}
			if got[0].severity != tc.want {
				t.Errorf("reported %s, want %s: %s", got[0].Severity, tc.want, got[0].Message)
			}
			if tc.contains != "" && !strings.Contains(got[0].Message, tc.contains) {
				t.Errorf("message %q does not contain %q", got[0].Message, tc.contains)
			}
			if got[0].File != "scripts/deploy.sh" || got[0].Line != 2 {
				t.Errorf("reported %s:%d, want scripts/deploy.sh:2", got[0].File, got[0].Line)
			}
		})
	}
}

func composeWith(service string) string {
	return "version: '3.8'\nservices:\n" + service
}

// The compose file is the other way this repository's container is stopped -
// `make run` starts it that way - and it fails identically: the default is ten
// seconds and it is nowhere near the budget it has to cover.
func TestComposeStopGraceAgainstTheShutdownBudget(t *testing.T) {
	for _, tc := range []struct {
		name     string
		service  string
		want     Severity
		contains string
	}{
		{
			name:    "generous",
			service: "  api:\n    image: go-admin:latest\n    stop_grace_period: 30s\n",
			want:    -1,
		},
		{
			name:     "not set at all",
			service:  "  api:\n    image: go-admin:latest\n",
			want:     Error,
			contains: "Set stop_grace_period: 13s",
		},
		{
			name:     "shorter than the shutdown",
			service:  "  api:\n    image: go-admin:latest\n    stop_grace_period: 5s\n",
			want:     Error,
			contains: "stop_grace_period on service api allows 5 seconds",
		},
		{
			name:    "longer than the shutdown but inside the margin",
			service: "  api:\n    image: go-admin:latest\n    stop_grace_period: 10s\n",
			want:    Warn,
		},
		{
			// Compose takes hours and minutes as well as seconds, and a check
			// that only read the digits would call 1m30s ninety times too
			// short.
			name:    "minutes and seconds",
			service: "  api:\n    image: go-admin:latest\n    stop_grace_period: 1m30s\n",
			want:    -1,
		},
		{
			// A service running something else is not this process, and its
			// grace period has nothing to do with this budget.
			name:    "another image is not this application",
			service: "  db:\n    image: mysql:8\n",
			want:    -1,
		},
		{
			// Built from this repository, so it is this application whatever
			// the image ends up being called.
			name:     "built here rather than named",
			service:  "  api:\n    build: .\n",
			want:     Error,
			contains: "service api, which sets no stop_grace_period",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{
				"config/extend.go":    hostConfigSource,
				"config/settings.yml": factorySettings,
				"docker-compose.yml":  composeWith(tc.service),
			})
			got := only(t, check(t, root, options{}), checkDockerStop)
			if tc.want < 0 {
				if len(got) != 0 {
					t.Fatalf("reported %d findings, want none:\n%v", len(got), got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("reported %d findings, want 1:\n%v", len(got), got)
			}
			if got[0].severity != tc.want {
				t.Errorf("reported %s, want %s: %s", got[0].Severity, tc.want, got[0].Message)
			}
			if tc.contains != "" && !strings.Contains(got[0].Message, tc.contains) {
				t.Errorf("message %q does not contain %q", got[0].Message, tc.contains)
			}
			if got[0].File != "docker-compose.yml" {
				t.Errorf("reported against %s, want docker-compose.yml", got[0].File)
			}
		})
	}
}

func TestComposeDurations(t *testing.T) {
	for _, tc := range []struct {
		in     string
		want   int
		wantOK bool
	}{
		{in: "30s", want: 30, wantOK: true},
		{in: "30", want: 30, wantOK: true},
		{in: "1m30s", want: 90, wantOK: true},
		{in: "2m", want: 120, wantOK: true},
		{in: "1h", want: 3600, wantOK: true},
		{in: "1h0m30s", want: 3630, wantOK: true},
		{in: "", wantOK: false},
		{in: "forever", wantOK: false},
		{in: "500ms", wantOK: false},
	} {
		t.Run(tc.in, func(t *testing.T) {
			got, ok := composeSeconds(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("composeSeconds(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
			}
			if ok && got != tc.want {
				t.Errorf("composeSeconds(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// The command can be written in a workflow or in the Makefile as easily as in
// a shell script, and a check that only looked at one of them would be quiet
// about the others.
func TestDockerStopIsFoundInEveryKindOfFile(t *testing.T) {
	root := fixture(t, map[string]string{
		"config/extend.go":           hostConfigSource,
		"config/settings.yml":        factorySettings,
		".github/workflows/ship.yml": "jobs:\n  deploy:\n    steps:\n      - run: docker stop app\n",
		"Makefile":                   "stop:\n\tdocker stop app\n",
		"scripts/deploy.sh":          "docker stop app\n",
	})
	got := only(t, check(t, root, options{}), checkDockerStop)
	if len(got) != 3 {
		t.Fatalf("found %d commands, want 3:\n%v", len(got), got)
	}
}

func TestLineOfIgnoresComments(t *testing.T) {
	content := []byte("a: 1\n  # terminationGracePeriodSeconds: 99\n  terminationGracePeriodSeconds: 30\n")
	if got := lineOf(content, "terminationGracePeriodSeconds:"); got != 3 {
		t.Errorf("lineOf = %d, want 3 - a commented-out key is not the setting", got)
	}
}
