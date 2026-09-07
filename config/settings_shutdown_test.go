package config

import (
	"testing"

	coreconfig "github.com/go-admin-team/go-admin-core/v2/config"
	"github.com/go-admin-team/go-admin-core/v2/config/source/file"
)

// shippedSettings is the shape the loader fills in, cut down to the part under
// test. The reader is JSON-based, so the keys are matched against field names
// case-insensitively - which is exactly the matching that silently drops a
// section the struct has no field for.
type shippedSettings struct {
	Settings struct {
		Extend Extend
	}
}

func (*shippedSettings) OnChange() {}

// The shutdown section has to arrive where it is read from, and with the
// values the documentation claims.
//
// This is the failure this batch exists to remove, one level up: the loader
// discards keys no field matches, without an error and without a log line, so
// a section put in the wrong place is written, accepted, and never applied.
// Nothing but loading the shipped file through the real loader can tell the
// two apart - the struct compiles either way.
//
// The values are asserted as well as the arrival. A settings.yml that shipped
// a different default from config.Default*Seconds would give two answers to
// "what does an unconfigured deployment spend", and the file is the one people
// read.
func TestTheShippedSettingsReachTheShutdownStruct(t *testing.T) {
	for _, name := range []string{"settings.yml", "settings.full.yml"} {
		t.Run(name, func(t *testing.T) {
			var loaded shippedSettings
			c, err := coreconfig.NewConfig(
				coreconfig.WithSource(file.NewSource(file.WithPath(name))),
				coreconfig.WithEntity(&loaded),
			)
			if err != nil {
				t.Fatalf("load %s: %v", name, err)
			}
			t.Cleanup(func() { _ = c.Close() })

			s := loaded.Settings.Extend.Shutdown
			if s.Drain == nil || s.Server == nil || s.Cleanup == nil {
				t.Fatalf("%s left extend.shutdown unfilled (%+v); the section is written but nothing reads it",
					name, s)
			}

			want := ShutdownBudget{
				Drain:   DefaultDrainSeconds,
				Server:  DefaultServerSeconds,
				Cleanup: DefaultCleanupSeconds,
			}
			got, err := s.Budget()
			if err != nil {
				t.Fatalf("%s does not resolve: %v", name, err)
			}
			if got != want {
				t.Errorf("%s ships %+v, want the documented defaults %+v", name, got, want)
			}
		})
	}
}
