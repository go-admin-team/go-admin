package config

import "testing"

func TestObjectStoreConfigured(t *testing.T) {
	if (ObjectStore{}).Configured() {
		t.Fatal("empty store reported as configured")
	}
	full := ObjectStore{Endpoint: "e", AccessKeyID: "a", AccessKeySecret: "s", BucketName: "b"}
	if !full.Configured() {
		t.Fatal("complete store reported as unconfigured")
	}
	if (ObjectStore{Endpoint: "e", AccessKeyID: "a"}).Configured() {
		t.Fatal("partial store reported as configured")
	}
}

func TestRateLimitThreshold(t *testing.T) {
	// Absent is the case an existing settings.yml hits after an upgrade: it has
	// no ratelimit section, and must keep the limit it always had.
	if got := (RateLimit{}).Threshold(); got != DefaultInboundQPS {
		t.Errorf("unconfigured limit = %v, want the default %v", got, DefaultInboundQPS)
	}

	zero := 0.0
	if got := (RateLimit{InboundQPS: &zero}).Threshold(); got != 0 {
		t.Errorf("explicit zero = %v, want 0 so the limiter can be turned off", got)
	}

	custom := 1500.0
	if got := (RateLimit{InboundQPS: &custom}).Threshold(); got != custom {
		t.Errorf("configured limit = %v, want %v", got, custom)
	}
}
