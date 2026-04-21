package schedule_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

// TestMultipleConsumers verifies that a single Scheduler drives a worker loop
// correctly when multiple goroutines coordinate via shared state.
func TestMultipleConsumers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s := schedule.New(100 * time.Millisecond)
	ch := s.Run(ctx)

	var total atomic.Int64
	done := make(chan struct{})

	go func() {
		defer close(done)
		for range ch {
			total.Add(1)
			if total.Load() >= 5 {
				cancel()
				return
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("worker loop did not complete in time")
	}

	if n := total.Load(); n < 5 {
		t.Fatalf("expected >= 5 ticks, got %d", n)
	}
}

// TestScheduler_WithJitter ensures the jitter option is accepted without panic.
func TestScheduler_WithJitter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	s := schedule.New(150*time.Millisecond, schedule.WithJitter(50*time.Millisecond))
	ch := s.Run(ctx)

	select {
	case _, ok := <-ch:
		if !ok {
			t.Fatal("channel closed immediately")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("no tick received with jitter option set")
	}
}
