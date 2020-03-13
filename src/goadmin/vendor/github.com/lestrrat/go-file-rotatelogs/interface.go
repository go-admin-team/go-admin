package rotatelogs

import (
	"os"
	"sync"
	"time"

	strftime "github.com/lestrrat/go-strftime"
)

// RotateLogs represents a log file that gets
// automatically rotated as you write to it.
type RotateLogs struct {
	clock         Clock
	curFn         string
	globPattern   string
	linkName      string
	maxAge        time.Duration
	mutex         sync.RWMutex
	outFh         *os.File
	pattern       *strftime.Strftime
	rotationTime  time.Duration
	rotationCount int
}

// Clock is the interface used by the RotateLogs
// object to determine the current time
type Clock interface {
	Now() time.Time
}
type clockFn func() time.Time

// UTC is an object satisfying the Clock interface, which
// returns the current time in UTC
var UTC = clockFn(func() time.Time { return time.Now().UTC() })

// Local is an object satisfying the Clock interface, which
// returns the current time in the local timezone
var Local = clockFn(time.Now)

// Option is used to pass optional arguments to
// the RotateLogs constructor
type Option interface {
	Configure(*RotateLogs) error
}

// OptionFn is a type of Option that is represented
// by a single function that gets called for Configure()
type OptionFn func(*RotateLogs) error
