package api

import (
	"strings"
	"sync"
	"testing"
	"time"

	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
)

// recordingQueue records what was done to it, in order. Register and Run are
// the two calls whose order is the point of this file; Append and Shutdown are
// here to satisfy the interface.
type recordingQueue struct {
	mu     sync.Mutex
	events []string
	ran    chan struct{}
}

func newRecordingQueue() *recordingQueue {
	return &recordingQueue{ran: make(chan struct{}, 4)}
}

func (q *recordingQueue) record(e string) {
	q.mu.Lock()
	q.events = append(q.events, e)
	q.mu.Unlock()
}

func (q *recordingQueue) seen() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]string(nil), q.events...)
}

func (q *recordingQueue) String() string                    { return "recording" }
func (q *recordingQueue) Append(corestorage.Messager) error { return nil }
func (q *recordingQueue) Register(name string, _ corestorage.ConsumerFunc) {
	q.record("register:" + name)
}
func (q *recordingQueue) Shutdown() {}

func (q *recordingQueue) Run() {
	q.record("run")
	select {
	case q.ran <- struct{}{}:
	default:
	}
}

// waitForRun waits for Run, which is started on a goroutine.
func (q *recordingQueue) waitForRun(t *testing.T) {
	t.Helper()
	select {
	case <-q.ran:
	case <-time.After(5 * time.Second):
		t.Fatalf("Run was never called; saw: %v", q.seen())
	}
}

// The consumers must be registered before the queue is started. A queue that
// is already running refuses further registration - the contract
// implementations answer storage.ErrQueueAlreadyStarted - and the legacy
// adapter this path goes through drops that error, so the wrong order loses
// consumers with nothing said about it. The memory backend does not care,
// which is exactly why this cannot be left to be noticed in use.
func TestConsumersAreRegisteredBeforeTheQueueIsStarted(t *testing.T) {
	attachedQueue.Store(0)
	t.Cleanup(func() { attachedQueue.Store(0) })

	q := newRecordingQueue()
	attachConsumersOnce(1, q)
	q.waitForRun(t)

	seen := q.seen()
	runAt := -1
	registers := 0
	for i, e := range seen {
		switch {
		case e == "run":
			if runAt < 0 {
				runAt = i
			}
		case strings.HasPrefix(e, "register:"):
			registers++
			if runAt >= 0 {
				t.Errorf("%q came after Run; a running queue refuses registration", e)
			}
		}
	}
	if registers != 3 {
		t.Errorf("registered %d consumers, want 3; saw %v", registers, seen)
	}
	if runAt < 0 {
		t.Errorf("the queue was never started; saw %v", seen)
	}
}

// AfterResource runs again on every configuration reload, so the hook has to
// be idempotent with respect to a given queue - not "does nothing the second
// time". Registering twice on the same queue would give every message two
// consumers and write every log row twice.
func TestTheSameQueueIsNotGivenConsumersTwice(t *testing.T) {
	attachedQueue.Store(0)
	t.Cleanup(func() { attachedQueue.Store(0) })

	q := newRecordingQueue()
	attachConsumersOnce(1, q)
	q.waitForRun(t)
	attachConsumersOnce(1, q)

	// Nothing to wait for on the second call, so give a wrong implementation
	// the time it would need to show up.
	time.Sleep(200 * time.Millisecond)
	if n := len(q.seen()); n != 4 {
		t.Errorf("%d calls after attaching twice to the same queue, want 4 (3 registers + 1 run); saw %v", n, q.seen())
	}
}

// The other half of the same rule: a reload builds a new adapter, and the
// consumers on the old one are attached to a queue nobody publishes to any
// more. A new generation must get its own set.
func TestANewQueueGetsItsOwnConsumers(t *testing.T) {
	attachedQueue.Store(0)
	t.Cleanup(func() { attachedQueue.Store(0) })

	first := newRecordingQueue()
	attachConsumersOnce(1, first)
	first.waitForRun(t)

	second := newRecordingQueue()
	attachConsumersOnce(2, second)
	second.waitForRun(t)

	if n := len(second.seen()); n != 4 {
		t.Errorf("the queue from the second generation saw %d calls, want 4; saw %v", n, second.seen())
	}
	if n := len(first.seen()); n != 4 {
		t.Errorf("the queue from the first generation saw %d calls, want 4 - it should not have been touched again; saw %v", n, first.seen())
	}
}

// Generation 0 means the configuration has no queue section at all, so nothing
// was installed and the runtime hands back its own memory queue. That case
// still has to get consumers - the registration it replaces was unconditional,
// and dropping it would stop the login and operation logs for anyone who
// commented the section out.
func TestAnUnconfiguredQueueStillGetsConsumers(t *testing.T) {
	attachedQueue.Store(0)
	t.Cleanup(func() { attachedQueue.Store(0) })

	q := newRecordingQueue()
	attachConsumersOnce(0, q)
	q.waitForRun(t)

	if n := len(q.seen()); n != 4 {
		t.Errorf("an unconfigured queue saw %d calls, want 4; saw %v", n, q.seen())
	}

	// And still only once.
	attachConsumersOnce(0, q)
	time.Sleep(200 * time.Millisecond)
	if n := len(q.seen()); n != 4 {
		t.Errorf("generation 0 was attached to twice: %d calls, want 4; saw %v", n, q.seen())
	}
}
