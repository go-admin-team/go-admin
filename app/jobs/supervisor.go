package jobs

import (
	"context"
	"sync"
	"time"

	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg/cronjob"
	"gorm.io/gorm"
)

// leaseTTL is how long a lease stays valid without being renewed, and
// leaseHeartbeat is how often the holder renews it.
//
// The gap between them is the point: at a third of the TTL, two consecutive
// renewals can fail - a restarting database, a paused container - and the
// third still lands before anything else may take the lease. Making them
// equal would hand the scheduler to another instance on the first missed
// beat.
//
// The TTL is also the longest the jobs can be stopped everywhere: an
// instance killed without running its shutdown leaves its lease behind, and
// the successor waits this long before taking it.
const (
	leaseTTL       = 30 * time.Second
	leaseHeartbeat = 10 * time.Second
)

// supervisor keeps one tenant's scheduler in step with one lease.
//
// It exists because holding the lease is not a decision made once at
// startup. An instance that never gets the lease has to keep asking, or the
// death of the current holder would stop the jobs until somebody restarted a
// process by hand; and an instance that holds it has to stop scheduling the
// moment it can no longer prove it still does, or a network partition turns
// into the two-schedulers-at-once defect (#915) that the lease exists to
// prevent.
type supervisor struct {
	key   string
	db    *gorm.DB
	lease *lease

	mu      sync.Mutex
	running bool
	// lastRenew is when this instance last proved it holds the lease. It
	// is compared only against this process's own later readings, never
	// against another instance's, so the monotonic clock is the right one
	// here - the reason the lease itself reads the database's clock does
	// not apply to measuring how long ago something happened locally.
	lastRenew time.Time

	stop     chan struct{}
	stopOnce sync.Once
	done     chan struct{}
}

func newSupervisor(key string, db *gorm.DB, owner string) *supervisor {
	return &supervisor{
		key:   key,
		db:    db,
		lease: &lease{db: db, owner: owner, ttl: leaseTTL},
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}
}

// start takes the lease if it is free, schedules this tenant's jobs if it
// got it, and then keeps both facts true for the life of the process.
//
// The first attempt is synchronous so that a single-instance deployment -
// which is nearly all of them - has its jobs registered by the time Setup
// returns, exactly as it did before there was a lease.
func (s *supervisor) start() {
	s.tick()

	go s.heartbeat()

	sdk.Runtime.SetShutdown(func(ctx context.Context) {
		s.shutdown(ctx)
	})
}

func (s *supervisor) heartbeat() {
	defer close(s.done)

	t := time.NewTicker(leaseHeartbeat)
	defer t.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-t.C:
			s.tick()
		}
	}
}

// tick asks for the lease and makes the scheduler match the answer.
func (s *supervisor) tick() {
	held, err := s.lease.acquire()
	if err != nil {
		// Not knowing is not the same as having lost it. The lease is
		// still ours until it expires, so the scheduler keeps running
		// and this instance keeps trying - a database that is briefly
		// unreachable must not stop the jobs, and must not hand them to
		// anyone else either, because nobody else can reach it to take
		// the lease.
		log.Errorf("[Job] scheduler lease for %s: %v", s.key, err)
		if s.heldFor() > leaseTTL {
			log.Errorf("[Job] scheduler lease for %s has not been renewed in %v; stopping the scheduler before anything else takes it",
				s.key, leaseTTL)
			s.stopScheduling()
		}
		return
	}

	if !held {
		s.stopScheduling()
		return
	}

	s.mu.Lock()
	s.lastRenew = time.Now()
	already := s.running
	s.mu.Unlock()

	if !already {
		s.startScheduling()
	}
}

// heldFor reports how long it has been since this instance last proved it
// holds the lease. A zero lastRenew means it never has, which is not a lease
// that has gone stale.
func (s *supervisor) heldFor() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastRenew.IsZero() {
		return 0
	}
	return time.Since(s.lastRenew)
}

func (s *supervisor) startScheduling() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	log.Infof("[Job] holding the scheduler lease for %s; registering its jobs", s.key)
	setup(s.key, s.db)
}

// stopScheduling stops this tenant's scheduler and puts a fresh one in its
// place.
//
// Fresh, rather than reusing the stopped one, because taking the lease back
// runs setup again and setup adds every enabled job to whatever scheduler is
// registered. Reusing it would leave the previous registration in place and
// fire every job twice - the symptom this whole change is here to remove,
// reintroduced one layer down.
func (s *supervisor) stopScheduling() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	log.Infof("[Job] no longer holding the scheduler lease for %s; stopping its jobs", s.key)
	ctx, cancel := context.WithTimeout(context.Background(), leaseHeartbeat)
	defer cancel()
	stopCrontab(ctx, s.key)
	sdk.Runtime.SetCrontabByTenant(s.key, cronjob.NewWithSeconds())
}

// shutdown stops the heartbeat, stops the scheduler and hands the lease back
// so a successor can take it now rather than waiting out the TTL.
func (s *supervisor) shutdown(ctx context.Context) {
	s.stopOnce.Do(func() { close(s.stop) })
	select {
	case <-s.done:
	case <-ctx.Done():
	}

	s.mu.Lock()
	wasRunning := s.running
	s.running = false
	s.mu.Unlock()

	if wasRunning {
		stopCrontab(ctx, s.key)
		if err := s.lease.release(); err != nil {
			log.Errorf("[Job] releasing the scheduler lease for %s: %v", s.key, err)
		}
	}
}

// holdsLease reports whether this instance is currently scheduling. It exists
// for the tests: everything else acts on the answer inside tick.
func (s *supervisor) holdsLease() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
