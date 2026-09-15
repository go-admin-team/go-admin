package api

import (
	"testing"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// The warning fires exactly where the generator's writing endpoints are served
// and nothing else refuses them.
//
// dev is the case the warning exists for: it is the shipped default, so it is
// the mode a deployment that changed nothing is running in. demo serves the
// routes too, but DemoEvn refuses all three by name, so warning there would
// describe an exposure that is not there.
func TestGeneratorWriteRoutesWarningFiresWhereTheExposureIs(t *testing.T) {
	for _, tc := range []struct {
		mode string
		want bool
		why  string
	}{
		{"dev", true, "shipped default, endpoints served and not refused"},
		{"demo", false, "served, but DemoEvn refuses all three"},
		{"test", false, "not served"},
		{"prod", false, "not served"},
		{"", false, "not served"},
	} {
		t.Run("mode="+tc.mode, func(t *testing.T) {
			previous := config.ApplicationConfig.Mode
			t.Cleanup(func() { config.ApplicationConfig.Mode = previous })
			config.ApplicationConfig.Mode = tc.mode

			if got := generatorWriteRoutesNeedWarning(); got != tc.want {
				t.Errorf("mode %q: warning = %v, want %v (%s)", tc.mode, got, tc.want, tc.why)
			}
		})
	}
}
