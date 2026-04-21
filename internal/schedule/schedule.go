// Package schedule provides interval-based tick scheduling for the poller.
package schedule

import (
	"context"
	"time"
)

// Scheduler emits ticks at a fixed interval until the context is cancelled.
type Scheduler struct {
	interval time.Duration
	jitter   time.Duration
}

// Option configures a Scheduler.
type Option func(*Scheduler)

// WithJitter adds up to jitter duration of random delay to each tick.
func WithJitter(j time.Duration) Option {
	return func(s *Scheduler) {
		if j > 0 {
			s.jitter = j
		}
	}
}

// New creates a Scheduler that ticks every interval.
func New(interval time.Duration, opts ...Option) *Scheduler {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	s := &Scheduler{interval: interval}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Run sends the current time on the returned channel at each tick.
// The channel is closed when ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) <-chan time.Time {
	ch := make(chan time.Time, 1)
	go func() {
		defer close(ch)
		// Fire immediately on start.
		select {
		case ch <- time.Now():
		case <-ctx.Done():
			return
		}
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				select {
				case ch <- t:
				default: // drop tick if consumer is slow
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

// Interval returns the configured tick interval.
func (s *Scheduler) Interval() time.Duration { return s.interval }
