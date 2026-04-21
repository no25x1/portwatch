// Package circuitbreaker implements a simple circuit breaker pattern
// for protecting downstream port scan targets from repeated failures.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Breaker is a circuit breaker that tracks failures for a given key.
type Breaker struct {
	mu           sync.Mutex
	threshold    int
	cooldown     time.Duration
	failures     map[string]int
	openUntil    map[string]time.Time
	now          func() time.Time
}

// New returns a Breaker that opens after threshold consecutive failures
// and resets after cooldown has elapsed.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		failures:  make(map[string]int),
		openUntil: make(map[string]time.Time),
		now:       time.Now,
	}
}

// Allow returns nil if the circuit is closed or half-open for key,
// or ErrOpen if the circuit is still open.
func (b *Breaker) Allow(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if until, ok := b.openUntil[key]; ok {
		if b.now().Before(until) {
			return ErrOpen
		}
		// cooldown elapsed — move to half-open
		delete(b.openUntil, key)
		b.failures[key] = 0
	}
	return nil
}

// RecordSuccess resets the failure count for key.
func (b *Breaker) RecordSuccess(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key] = 0
	delete(b.openUntil, key)
}

// RecordFailure increments the failure count and opens the circuit
// if the threshold has been reached.
func (b *Breaker) RecordFailure(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key]++
	if b.failures[key] >= b.threshold {
		b.openUntil[key] = b.now().Add(b.cooldown)
	}
}

// StateOf returns the current State for key.
func (b *Breaker) StateOf(key string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	until, ok := b.openUntil[key]
	if !ok {
		return StateClosed
	}
	if b.now().Before(until) {
		return StateOpen
	}
	return StateHalfOpen
}
