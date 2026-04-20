package throttle_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/throttle"
)

// TestConcurrentAcquire verifies that concurrent goroutines never exceed the
// token budget within a single window.
func TestConcurrentAcquire_NeverExceedsBudget(t *testing.T) {
	const budget = 10
	th := throttle.New(budget, time.Minute)

	var acquired atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < budget; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			if err := th.Acquire(ctx); err == nil {
				acquired.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := acquired.Load(); got != budget {
		t.Fatalf("expected %d acquisitions, got %d", budget, got)
	}
	if th.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %d", th.Remaining())
	}
}

// TestThrottle_ZeroWindowStillRefills ensures a very short window causes
// rapid refills and does not deadlock.
func TestThrottle_ShortWindowAllowsBurst(t *testing.T) {
	th := throttle.New(2, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var count int
	for count < 6 {
		if err := th.Acquire(ctx); err != nil {
			t.Fatalf("unexpected error after %d acquires: %v", count, err)
		}
		count++
	}
}
