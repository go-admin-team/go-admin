package jobs

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

// The frame the leaked goroutines park in. Remove starts it, and it is the
// only goroutine in this package that sends on a channel the caller may have
// walked away from.
const removeSenderFrame = "go-admin/app/jobs.Remove.func1"

// Remove hands the caller a channel it is free to abandon: RemoveJob stops
// waiting after a second and returns a timeout error. The send therefore has
// to complete with nobody receiving, or every stop that times out parks a
// goroutine on it for the life of the process.
//
// The order matters. Counting parked goroutines straight after calling Remove
// would pass while proving nothing, because the goroutine may not have reached
// the send yet. So the entries are waited out first: an empty scheduler means
// every goroutine is at or past its send, and only then is a survivor a leak.
func TestRemoveLetsItsGoroutineFinishWithNobodyReceiving(t *testing.T) {
	const jobs = 20

	c := cron.New()
	ids := make([]cron.EntryID, 0, jobs)
	for i := 0; i < jobs; i++ {
		id, err := c.AddFunc("@every 1h", func() {})
		if err != nil {
			t.Fatalf("AddFunc: %v", err)
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		// The returned channel is dropped on purpose: this is what a caller
		// that has already timed out leaves behind.
		_ = Remove(c, int(id))
	}

	if err := waitFor(3*time.Second, func() bool { return len(c.Entries()) == 0 }); err != nil {
		t.Fatalf("the scheduler still holds %d entries, so the goroutines never reached their send "+
			"and this test cannot show anything", len(c.Entries()))
	}

	if err := waitFor(3*time.Second, func() bool { return parkedInRemove() == 0 }); err != nil {
		t.Errorf("%d of %d goroutines are still parked sending on an abandoned channel:\n%s",
			parkedInRemove(), jobs, oneParkedStack())
	}
}

func waitFor(d time.Duration, done func() bool) error {
	deadline := time.Now().Add(d)
	for {
		if done() {
			return nil
		}
		if time.Now().After(deadline) {
			return errTimeout
		}
		time.Sleep(10 * time.Millisecond)
	}
}

var errTimeout = timeoutError{}

type timeoutError struct{}

func (timeoutError) Error() string { return "timed out" }

func parkedInRemove() int {
	return strings.Count(goroutineDump(), removeSenderFrame)
}

func oneParkedStack() string {
	for _, block := range strings.Split(goroutineDump(), "\n\n") {
		if strings.Contains(block, removeSenderFrame) {
			return block
		}
	}
	return "(none)"
}

func goroutineDump() string {
	buf := make([]byte, 1<<20)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
