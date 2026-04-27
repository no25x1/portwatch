// Package semaphore provides a weighted semaphore for bounding concurrent
// port-scan goroutines across multiple hosts.
package semaphore

import (
	"context"
	"fmt"
)

// Semaphore is a counting semaphore backed by a buffered channel.
type Semaphore struct {
	ch chan struct{}
}

// New returns a Semaphore with the given capacity.
// capacity must be greater than zero.
func New(capacity int) (*Semaphore, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("semaphore: capacity must be > 0, got %d", capacity)
	}
	return &Semaphore{ch: make(chan struct{}, capacity)}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a slot is obtained.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees one slot. It panics if called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Available returns the number of slots that can be acquired without blocking.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}

// Cap returns the total capacity of the semaphore.
func (s *Semaphore) Cap() int {
	return cap(s.ch)
}
