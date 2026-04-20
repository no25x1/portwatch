// Package ratelimit provides a simple per-key rate limiter to suppress
// repeated alerts for the same host:port within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per key and suppresses duplicates
// that arrive within the cooldown duration.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Limiter with the given cooldown window.
// A cooldown of zero disables rate limiting (every call is allowed).
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether the event identified by key should be allowed through.
// It returns true the first time a key is seen, and again only after the
// cooldown window has elapsed since the previous allowed event.
func (l *Limiter) Allow(key string) bool {
	if l.cooldown == 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, causing the next call to
// Allow for that key to be permitted immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Len returns the number of keys currently tracked.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.last)
}
