package limiter_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/limiter"
)

func TestLimiter_ScanSimulation(t *testing.T) {
	const concurrency = 4
	const tasks = 40

	l, err := limiter.New(concurrency)
	if err != nil {
		t.Fatal(err)
	}

	var active int64
	var peak int64
	var wg sync.WaitGroup

	for i := 0; i < tasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(context.Background()); err != nil {
				t.Errorf("Acquire: %v", err)
				return
			}
			cur := atomic.AddInt64(&active, 1)
			for {
				old := atomic.LoadInt64(&peak)
				if cur <= old || atomic.CompareAndSwapInt64(&peak, old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt64(&active, -1)
			l.Release()
		}()
	}
	wg.Wait()

	if peak > int64(concurrency) {
		t.Fatalf("peak concurrent=%d exceeded limit=%d", peak, concurrency)
	}
	if l.Active() != 0 {
		t.Fatalf("expected 0 active after all goroutines finished, got %d", l.Active())
	}
}
