package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/retry"
)

var errTemp = errors.New("temporary failure")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	p := retry.DefaultPolicy()
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1}
	var calls int32
	err := p.Do(context.Background(), func() error {
		if atomic.AddInt32(&calls, 1) < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1}
	var calls int32
	err := p.Do(context.Background(), func() error {
		atomic.AddInt32(&calls, 1)
		return errTemp
	})
	if !errors.Is(err, retry.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	p := retry.Policy{MaxAttempts: 5, InitialDelay: 50 * time.Millisecond, Multiplier: 1}
	ctx, cancel := context.WithCancel(context.Background())
	var calls int32
	err := p.Do(ctx, func() error {
		if atomic.AddInt32(&calls, 1) == 1 {
			cancel()
		}
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	p := retry.Policy{MaxAttempts: 0, InitialDelay: time.Millisecond}
	calls := 0
	_ = p.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
