// Package dedupe provides event deduplication to suppress repeated
// notifications for ports that have not changed state since the last alert.
package dedupe

import (
	"sync"
	"time"
)

// StateKey uniquely identifies a host/port pair.
type StateKey struct {
	Host string
	Port int
}

// entry records the last-alerted state for a key.
type entry struct {
	open      bool
	alertedAt time.Time
}

// Filter suppresses duplicate events: an event is only passed through when
// the open/closed state differs from the previously alerted state, or when
// the cooldown window has elapsed.
type Filter struct {
	mu       sync.Mutex
	seen     map[StateKey]entry
	cooldown time.Duration
	now      func() time.Time
}

// New creates a Filter with the given cooldown duration.
// A zero cooldown means events are only suppressed when state is identical
// and no time constraint applies.
func New(cooldown time.Duration) *Filter {
	return &Filter{
		seen:     make(map[StateKey]entry),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true if the event for (host, port, open) should be forwarded.
// It updates internal state when the event is allowed through.
func (f *Filter) Allow(host string, port int, open bool) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := StateKey{Host: host, Port: port}
	now := f.now()

	if prev, ok := f.seen[key]; ok {
		// Same state and within cooldown — suppress.
		if prev.open == open && (f.cooldown == 0 || now.Sub(prev.alertedAt) < f.cooldown) {
			return false
		}
	}

	f.seen[key] = entry{open: open, alertedAt: now}
	return true
}

// Reset clears all tracked state, forcing the next event for every key
// to be forwarded regardless of previous history.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[StateKey]entry)
}

// Len returns the number of distinct host/port pairs currently tracked.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}
