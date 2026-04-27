package semaphore_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/portwatch/internal/semaphore"
)

// TestScanSimulation verifies that a pool of "scan" goroutines never exceeds
// the semaphore capacity, mirroring real portwatch scan fan-out.
func TestScanSimulation(t *testing.T) {
	const concurrency = 8
	const targets = 64

	sem, err := semaphore.New(concurrency)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var active atomic.Int64
	var violations atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < targets; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sem.Acquire(context.Background()); err != nil {
				return
			}
			defer sem.Release()
			v := active.Add(1)
			if v > concurrency {
				violations.Add(1)
			}
			time.Sleep(time.Millisecond)
			active.Add(-1)
		}()
	}
	wg.Wait()

	if v := violations.Load(); v > 0 {
		t.Errorf("%d concurrency violations detected", v)
	}
	if sem.Available() != concurrency {
		t.Errorf("expected all slots returned: available=%d cap=%d",
			sem.Available(), sem.Cap())
	}
}

// TestSemaphore_EarlyCancel verifies that cancelling a context unblocks
// waiting goroutines cleanly without deadlock.
func TestSemaphore_EarlyCancel(t *testing.T) {
	sem, _ := semaphore.New(1)
	_ = sem.Acquire(context.Background()) // fill the only slot

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := sem.Acquire(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("unblock took too long: %v", elapsed)
	}
}
