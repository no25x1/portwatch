// Package retry provides a simple exponential-backoff retry helper used
// by the poller and watcher to tolerate transient scan failures.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy describes how retries should be attempted.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	// A value of 1 means no retries.
	MaxAttempts int

	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration

	// Multiplier is applied to the delay after each failed attempt.
	// Values <= 1 are treated as 1 (constant backoff).
	Multiplier float64

	// MaxDelay caps the computed delay.
	MaxDelay time.Duration
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		Multiplier:   2.0,
		MaxDelay:     5 * time.Second,
	}
}

// ErrExhausted is returned when all attempts have been consumed.
var ErrExhausted = errors.New("retry: all attempts exhausted")

// Do calls fn up to p.MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
// The last non-nil error from fn is wrapped with ErrExhausted.
func (p Policy) Do(ctx context.Context, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	mult := p.Multiplier
	if mult <= 1 {
		mult = 1
	}

	delay := p.InitialDelay
	var lastErr error

	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * mult)
			if p.MaxDelay > 0 && delay > p.MaxDelay {
				delay = p.MaxDelay
			}
		}
	}

	return errors.Join(ErrExhausted, lastErr)
}
