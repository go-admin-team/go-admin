package file_store

import "fmt"

// ErrNotConfigured is returned by an object store that was never given
// credentials. The Client field of every implementation is an interface{}
// assigned in Setup, so an unconfigured store holds nil, and asserting nil to
// the provider's client type panics. The upload endpoint reaches this path
// whenever a request asks for a provider the deployment has not configured.
type ErrNotConfigured struct {
	Driver DriverType
}

func (e *ErrNotConfigured) Error() string {
	return fmt.Sprintf("file store %s is not configured; set it under extend.fileStore", e.Driver)
}

func notConfigured(driver DriverType) error {
	return &ErrNotConfigured{Driver: driver}
}
