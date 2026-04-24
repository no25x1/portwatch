package limiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAcquireRelease_BasicFlow(t *testing.T) {
	l, _ := New(2)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	if l.Active() != 1 {
		t.Fatalf("expected active=1, got %d", l.Active())
	}
	l.Release()
	if l.Active() != 0 {
		t.Fatalf("expected active=0 after release, got %d", l.Active())
	}
}

func TestAcquire_BlocksAtMax(t *testing.T) {
	l, _ := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := l.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer l.Release()

	err := l.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context deadline exceeded")
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	l, _ := New(1)
	_ = l.Acquire(context.Background()) // fill the slot
	defer l.Release()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := l.Acquire(ctx); err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestConcurrent_NeverExceedsMax(t *testing.T) {
	const max = 3
	const goroutines = 20
	l, _ := New(max)

	var mu sync.Mutex
	peak := 0
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(context.Background())
			mu.Lock()
			if a := l.Active(); a > peak {
				peak = a
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			l.Release()
		}()
	}
	wg.Wait()

	if peak > max {
		t.Fatalf("peak active %d exceeded max %d", peak, max)
	}
}

func TestRelease_PanicsWithoutAcquire(t *testing.T) {
	l, _ := New(2)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on Release without Acquire")
		}
	}()
	l.Release()
}
