package throttle

import (
	"context"
	"testing"
	"time"
)

func TestAcquire_UnlimitedAlwaysSucceeds(t *testing.T) {
	th := New(0, time.Second)
	for i := 0; i < 100; i++ {
		if err := th.Acquire(context.Background()); err != nil {
			t.Fatalf("unexpected error on unlimited throttle: %v", err)
		}
	}
}

func TestAcquire_ConsumesTokens(t *testing.T) {
	th := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := th.Acquire(context.Background()); err != nil {
			t.Fatalf("acquire %d failed: %v", i, err)
		}
	}
	if got := th.Remaining(); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}

func TestAcquire_BlocksWhenExhausted(t *testing.T) {
	th := New(1, time.Minute)
	_ = th.Acquire(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := th.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context deadline exceeded, got nil")
	}
}

func TestAcquire_RefillsAfterWindow(t *testing.T) {
	now := time.Now()
	th := New(2, 100*time.Millisecond)
	th.now = func() time.Time { return now }

	_ = th.Acquire(context.Background())
	_ = th.Acquire(context.Background())

	if th.Remaining() != 0 {
		t.Fatal("expected 0 remaining after exhaustion")
	}

	// advance clock past the window
	th.now = func() time.Time { return now.Add(200 * time.Millisecond) }

	if err := th.Acquire(context.Background()); err != nil {
		t.Fatalf("expected token after refill, got: %v", err)
	}
	if th.Remaining() != 1 {
		t.Fatalf("expected 1 remaining after refill and one acquire, got %d", th.Remaining())
	}
}

func TestAcquire_ContextCancelledImmediately(t *testing.T) {
	th := New(1, time.Minute)
	_ = th.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := th.Acquire(ctx); err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestRemaining_UnlimitedReturnsNegOne(t *testing.T) {
	th := New(0, time.Second)
	if got := th.Remaining(); got != -1 {
		t.Fatalf("expected -1 for unlimited, got %d", got)
	}
}
