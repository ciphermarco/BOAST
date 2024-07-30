package storage

import (
	"container/heap"
	"errors"
	"fmt"
	"hash"
	"sync"
	"time"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"

	"golang.org/x/crypto/blake2b"
)

// Config represents the storage's configurable options.
type Config struct {
	TTL             time.Duration
	CheckInterval   time.Duration
	MaxRestarts     int
	MaxEvents       int
	MaxEventsByTest int
	MaxDumpSize     int
	HMACKey         []byte
}

// Storage represents the storage itself, holding its configurations and state.
type Storage struct {
	mu          sync.RWMutex
	tests       map[string]test
	maxTests    int
	totalTests  int
	totalEvents int
	hmac        hash.Hash
	cfg         Config
}

// test represents a test of this application.
// A test is identified by an id. It holds a canary token to be used in the response
// from receivers when it may aid testing and recorded events for this test's id.
type test struct {
	id     string
	canary string
	events *eventHeap
}

// New contains the logic to construct and return a new *Storage according to the passed
// *Config object. In case of error, it returns the error to the caller.
func New(cfg *Config) (*Storage, error) {
	hmac, err := blake2b.New256(cfg.HMACKey)
	if err != nil {
		return nil, err
	}
	maxTests := 0
	if cfg.MaxEvents > 0 && cfg.MaxEventsByTest > 0 {
		maxTests = cfg.MaxEvents / cfg.MaxEventsByTest
	}
	s := &Storage{
		tests:    make(map[string]test),
		maxTests: maxTests,
		hmac:     hmac,
		cfg:      *cfg,
	}
	return s, nil
}

// SetTest creates a new test or fetches an existing one to return a newly generated or
// already existing test id. In case of error, it returns the error to the caller.
func (s *Storage) SetTest(secret []byte) (id string, canary string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sum := s.unsafeHmac(secret)
	id, canary = app.ToBase32(sum[:len(sum)/2]), app.ToBase32(sum[len(sum)/2:])
	if t, exists := s.tests[id]; exists {
		return t.id, t.canary, nil
	} else if s.totalTests < s.maxTests {
		events := &eventHeap{}
		heap.Init(events)
		s.tests[id] = test{
			id:     id,
			canary: canary,
			events: events,
		}
		s.totalTests++
		return id, canary, nil
	}
	return "", "", errors.New("could not create test")
}

// SearchTest receives a function to be run against each tests' id and canary, and
// returns the id and canary of the first test for which the passed function returns
// true. If the function never returns true empty strings are returned to the caller.
func (s *Storage) SearchTest(f func(k, v string) bool) (id string, canary string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for id, t := range s.tests {
		if f(id, t.canary) {
			return id, t.canary
		}
	}
	return "", ""
}

// StoreEvent appends an event to an existing test if it exists, otherwise it will
// return an error to the caller.
func (s *Storage) StoreEvent(evt app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := evt.TestID
	if t, exists := s.tests[id]; exists {
		if s.cfg.MaxEvents > 0 && s.cfg.MaxEventsByTest > 0 && s.totalEvents <= s.cfg.MaxEvents {
			if t.events.Len() >= s.cfg.MaxEventsByTest {
				s.unsafePopEvent(id)
			}
			if len(evt.Dump) > s.cfg.MaxDumpSize {
				evt.Dump = evt.Dump[:s.cfg.MaxDumpSize]
			}
			s.unsafePushEvent(id, evt)
		}
		return nil
	}
	return fmt.Errorf("test id %s does not exist", id)
}

// LoadEvents returns the copy of an test's events slice if the test exists.
func (s *Storage) LoadEvents(id string) (evts []app.Event, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if t, exists := s.tests[id]; exists {
		evts := make([]app.Event, t.events.Len())
		copy(evts, *t.events)
		return evts, true
	}
	return evts, false
}

// TotalTests returns the number of total tests recorded in the storage at the moment.
func (s *Storage) TotalTests() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalTests
}

// TotalEvents returns the number of total events recorded in the storage at the moment.
func (s *Storage) TotalEvents() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalEvents
}

// expire takes care of expiring (i.e. deleting) events according to the configured TTL
// and check interval.
func (s *Storage) expire() (err error) {
	defer func() error {
		if r := recover(); r != nil {
			err = fmt.Errorf("storage expiration error (panic): %v", r)
		}
		return err
	}()
	for range time.Tick(s.cfg.CheckInterval) {
		s.mu.RLock()
		for id, t := range s.tests {
			if t.events.Len() == 0 {
				s.mu.RUnlock()
				s.mu.Lock()

				s.unsafeDeleteTest(id)

				s.mu.Unlock()
				s.mu.RLock()
				continue
			}
			ttl := s.cfg.TTL
			for t.events.Len() > 0 && time.Since((*t.events)[0].Time) > ttl {
				s.mu.RUnlock()
				s.mu.Lock()

				s.unsafePopEvent(id)

				s.mu.Unlock()
				s.mu.RLock()
			}
		}
		s.mu.RUnlock()
	}
	return errors.New("storage expiration error")
}

// StartExpire is used by the caller to start expiring events and, in case of a panic
// from expire function, try to restart the expiration process the configured number of
// times.
//
// Panics should not and are not expected to happen as a normal occurrence, so this may
// be useless and soon to be dropped.
func (s *Storage) StartExpire(ret chan error) {
	err := s.expire()
	for i := 0; i < s.cfg.MaxRestarts; i++ {
		log.Info("Events expiration stopped. Restarting. (%d)\n", i+1)
		log.Debug("Storage.StartExpire error: %v", err)
		err = s.expire()
	}
	ret <- err
}

// unsafePushEvent pushes an event to an test's events leaving the mutex lock to the caller.
// It's unsafe to be used without setting the appropriate lock externally.
func (s *Storage) unsafePushEvent(id string, evt app.Event) {
	if t, exists := s.tests[id]; exists {
		heap.Push(t.events, evt)
		s.totalEvents++
	}

}

// unsafePopEvent pops an event from an test's events leaving the mutex lock to the caller.
// It's unsafe to be used without setting the appropriate lock externally.
func (s *Storage) unsafePopEvent(id string) {
	if t, exists := s.tests[id]; exists {
		heap.Pop(t.events)
		s.totalEvents--
		if t.events.Len() == 0 {
			s.unsafeDeleteTest(id)
		}
	}
}

func (s *Storage) unsafeDeleteTest(id string) {
	delete(s.tests, id)
	s.totalTests--
}

// unsafeHmac uses the storage's hmac and passed bytes to return an HMAC'd sum.
// It's unsafe to be used without setting the appropriate lock externally.
func (s *Storage) unsafeHmac(secret []byte) []byte {
	s.hmac.Reset()
	s.hmac.Write(secret)
	return s.hmac.Sum(nil)
}
