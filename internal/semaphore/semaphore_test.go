package semaphore

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_InvalidCapacity(t *testing.T) {
	for _, cap := range []int{0, -1, -100} {
		_, err := New(cap)
		if err == nil {
			t.Errorf("expected error for capacity %d", cap)
		}
	}
}

func TestNew_ValidCapacity(t *testing.T) {
	sem, err := New(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sem.Cap() != 4 {
		t.Fatalf("expected cap 4, got %d", sem.Cap())
	}
	if sem.Available() != 4 {
		t.Fatalf("expected available 4, got %d", sem.Available())
	}
}

func TestAcquireRelease_Basic(t *testing.T) {
	sem, _ := New(2)
	ctx := context.Background()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	if sem.Available() != 1 {
		t.Fatalf("expected 1 available, got %d", sem.Available())
	}
	sem.Release()
	if sem.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", sem.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	sem, _ := New(1)
	ctx := context.Background()
	_ = sem.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := sem.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	sem, _ := New(1)
	_ = sem.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := sem.Acquire(ctx); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestRelease_PanicsWithoutAcquire(t *testing.T) {
	sem, _ := New(2)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	sem.Release()
}

func TestConcurrent_NeverExceedsCap(t *testing.T) {
	const cap = 5
	const workers = 50
	sem, _ := New(cap)

	var active atomic.Int64
	var peak atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sem.Acquire(context.Background())
			v := active.Add(1)
			for {
				old := peak.Load()
				if v <= old || peak.CompareAndSwap(old, v) {
					break
				}
			}
			time.Sleep(2 * time.Millisecond)
			active.Add(-1)
			sem.Release()
		}()
	}
	wg.Wait()
	if p := peak.Load(); p > cap {
		t.Fatalf("peak concurrent %d exceeded cap %d", p, cap)
	}
}
