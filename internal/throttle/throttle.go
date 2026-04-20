// Package throttle provides a token-bucket style scan throttler that limits
// how many port checks can be dispatched per second across all hosts.
package throttle

import (
	"context"
	"sync"
	"time"
)

// Throttle controls the rate at which scan tokens are issued.
type Throttle struct {
	mu       sync.Mutex
	tokens   int
	max      int
	refillAt time.Time
	window   time.Duration
	now      func() time.Time
}

// New creates a Throttle that allows up to maxPerWindow scans within each
// window duration. A zero or negative maxPerWindow is treated as unlimited.
func New(maxPerWindow int, window time.Duration) *Throttle {
	return &Throttle{
		tokens: maxPerWindow,
		max:    maxPerWindow,
		window: window,
		now:    time.Now,
	}
}

// Acquire blocks until a scan token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is obtained.
func (t *Throttle) Acquire(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if t.tryAcquire() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// tryAcquire attempts to consume one token, refilling the bucket if the
// current window has elapsed. Returns true when a token was consumed.
func (t *Throttle) tryAcquire() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.max <= 0 {
		return true
	}

	now := t.now()
	if now.After(t.refillAt) {
		t.tokens = t.max
		t.refillAt = now.Add(t.window)
	}

	if t.tokens <= 0 {
		return false
	}
	t.tokens--
	return true
}

// Remaining returns the number of tokens left in the current window.
func (t *Throttle) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.max <= 0 {
		return -1
	}
	return t.tokens
}
