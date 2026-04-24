// Package limiter provides a concurrency limiter that caps the number of
// simultaneous port scans across all hosts, preventing resource exhaustion
// during large sweeps.
package limiter

import (
	"context"
	"fmt"
	"sync"
)

// Limiter controls the maximum number of concurrent operations.
type Limiter struct {
	mu      sync.Mutex
	sem     chan struct{}
	max     int
	active  int
}

// New creates a Limiter that allows at most max concurrent acquisitions.
// max must be greater than zero.
func New(max int) (*Limiter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("limiter: max must be > 0, got %d", max)
	}
	return &Limiter{
		sem: make(chan struct{}, max),
		max: max,
	}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if the context expires before a slot is obtained.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		l.mu.Lock()
		l.active++
		l.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot. Panics if called more times
// than Acquire has succeeded.
func (l *Limiter) Release() {
	select {
	case <-l.sem:
		l.mu.Lock()
		l.active--
		l.mu.Unlock()
	default:
		panic("limiter: Release called without matching Acquire")
	}
}

// Active returns the number of currently held slots.
func (l *Limiter) Active() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
}

// Max returns the configured concurrency ceiling.
func (l *Limiter) Max() int {
	return l.max
}
