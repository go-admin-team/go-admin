package storage

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
	"github.com/go-admin-team/go-admin-core/v2/storage/queue"
)

// redisAddrEnv points these tests at a server. They are skipped without it, so
// a developer with no redis running still gets a green run - and CI sets it,
// which is the point: the ordering rule they cover is invisible on the memory
// backend, and memory is the default. A suite that only ever exercised the
// default would report success for a queue that silently drops every consumer.
const redisAddrEnv = "GO_ADMIN_TEST_REDIS_ADDR"

func redisAddr(t *testing.T) string {
	t.Helper()
	addr := os.Getenv(redisAddrEnv)
	if addr != "" {
		return addr
	}
	// Skipping locally is the point; skipping in CI is the failure this whole
	// file exists to prevent. A workflow that renamed the variable, or dropped
	// the service, would otherwise go green while these two tests quietly did
	// nothing - which is the same shape as the defect they cover.
	if os.Getenv("CI") != "" {
		t.Fatalf("%s is not set while CI is: the redis-backed queue tests must not skip here", redisAddrEnv)
	}
	t.Skipf("%s is not set; skipping the redis-backed queue tests", redisAddrEnv)
	return ""
}

// newRedisQueue builds the queue the same way setupQueue does - through
// config.QueueConfig.Setup - so that what is under test is the adapter this
// repository actually gets, LegacyQueueAdapter and all, rather than a redis
// client wired up by the test.
func newRedisQueue(t *testing.T, prefix string) corestorage.AdapterQueue {
	t.Helper()
	previous := config.QueueConfig
	t.Cleanup(func() { config.QueueConfig = previous })

	config.QueueConfig = &config.Queue{
		Redis: &config.RedisQueue{
			RedisOptions: config.RedisOptions{Addr: redisAddr(t)},
			Group:        prefix,
			KeyPrefix:    prefix,
		},
	}
	q, err := config.QueueConfig.Setup()
	if err != nil {
		t.Fatalf("queue setup: %v", err)
	}
	t.Cleanup(q.Shutdown)
	return q
}

func message(t *testing.T, stream string) corestorage.Messager {
	t.Helper()
	m := &queue.Message{}
	m.SetStream(stream)
	m.SetValues(map[string]interface{}{"hello": "world"})
	return m
}

// Registered first, then started: the consumer gets the message. This is the
// order setupQueue and attachQueueConsumers now produce between them.
func TestRedisQueueDeliversToAConsumerRegisteredBeforeTheStart(t *testing.T) {
	stream := "t-ordered"
	q := newRedisQueue(t, "gotest-ordered")

	got := make(chan struct{}, 1)
	q.Register(stream, func(corestorage.Messager) error {
		select {
		case got <- struct{}{}:
		default:
		}
		return nil
	})
	go q.Run()

	// Give Start a moment to reach its read loop before publishing.
	time.Sleep(500 * time.Millisecond)
	if err := q.Append(message(t, stream)); err != nil {
		t.Fatalf("append: %v", err)
	}

	select {
	case <-got:
	case <-time.After(15 * time.Second):
		t.Fatal("the consumer never received the message")
	}
}

// Started first, then registered: the registration is refused and every
// publish afterwards fails.
//
// Subscribe answers ErrQueueAlreadyStarted, and LegacyQueueAdapter.Register
// returns nothing, so the caller cannot know - that part is silent. What is not
// silent is the consequence: no consumer group was created, so Publish refuses
// the topic with ErrNoHandler on every single request, and go-admin's call
// sites log that at error level while the login and operation log rows are
// never written.
//
// This is the test the memory backend cannot provide. queue.Memory's Register
// starts another consumer goroutine whatever the state, so the same code passes
// there - which is how the defect survived, memory being the default.
func TestRedisQueueRefusesAConsumerRegisteredAfterTheStart(t *testing.T) {
	stream := "t-late"
	q := newRedisQueue(t, "gotest-late")

	go q.Run()
	time.Sleep(500 * time.Millisecond)

	q.Register(stream, func(corestorage.Messager) error { return nil })

	err := q.Append(message(t, stream))
	if err == nil {
		t.Fatal("a message was accepted for a topic whose registration came after Start; " +
			"if the backend now accepts late registration, the ordering rule in setupQueue can be revisited")
	}
	if !errors.Is(err, corestorage.ErrNoHandler) {
		t.Fatalf("append failed with %v, want %v - the test is meant to pin the "+
			"missing-consumer path, not any error at all", err, corestorage.ErrNoHandler)
	}
}
