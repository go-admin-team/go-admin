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
