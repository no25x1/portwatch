package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

const key = "host:8080"

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(key); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		b.RecordFailure(key)
	}
	if err := b.Allow(key); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_RemainsClosedBelowThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure(key)
	b.RecordFailure(key)
	if err := b.Allow(key); err != nil {
		t.Fatalf("expected nil below threshold, got %v", err)
	}
}

func TestRecordSuccess_ResetsFailures(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure(key)
	b.RecordFailure(key)
	b.RecordSuccess(key)
	b.RecordFailure(key)
	b.RecordFailure(key)
	// only 2 failures after reset — should still be closed
	if err := b.Allow(key); err != nil {
		t.Fatalf("expected nil after partial reset, got %v", err)
	}
}

func TestStateOf_Transitions(t *testing.T) {
	now := time.Now()
	b := circuitbreaker.New(1, 50*time.Millisecond)
	// inject fake clock
	b2 := circuitbreaker.New(1, 50*time.Millisecond)
	_ = now

	if s := b2.StateOf(key); s != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed, got %v", s)
	}
	b2.RecordFailure(key)
	if s := b2.StateOf(key); s != circuitbreaker.StateOpen {
		t.Fatalf("expected Open, got %v", s)
	}
	time.Sleep(60 * time.Millisecond)
	if s := b2.StateOf(key); s != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", s)
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b := circuitbreaker.New(1, 40*time.Millisecond)
	b.RecordFailure(key)
	if err := b.Allow(key); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen immediately after open, got %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if err := b.Allow(key); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
}

func TestAllow_MultipleKeys_Independent(t *testing.T) {
	b := circuitbreaker.New(2, time.Minute)
	b.RecordFailure("a:80")
	b.RecordFailure("a:80")
	if err := b.Allow("b:80"); err != nil {
		t.Fatalf("unrelated key should remain closed, got %v", err)
	}
	if err := b.Allow("a:80"); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen for a:80, got %v", err)
	}
}
