package backoff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestDelay_BaseOnFirstAttempt(t *testing.T) {
	p := backoff.Policy{Base: 100 * time.Millisecond, Max: 10 * time.Second, Factor: 2.0}
	d := p.Delay(0)
	if d != 100*time.Millisecond {
		t.Fatalf("attempt 0: want 100ms, got %v", d)
	}
}

func TestDelay_GrowsExponentially(t *testing.T) {
	p := backoff.Policy{Base: 100 * time.Millisecond, Max: 10 * time.Second, Factor: 2.0}
	prev := p.Delay(0)
	for attempt := 1; attempt <= 4; attempt++ {
		curr := p.Delay(attempt)
		if curr <= prev {
			t.Fatalf("attempt %d: delay did not grow (prev=%v curr=%v)", attempt, prev, curr)
		}
		prev = curr
	}
}

func TestDelay_CappedAtMax(t *testing.T) {
	p := backoff.Policy{Base: 1 * time.Second, Max: 5 * time.Second, Factor: 10.0}
	for attempt := 0; attempt < 10; attempt++ {
		d := p.Delay(attempt)
		if d > p.Max {
			t.Fatalf("attempt %d: delay %v exceeds max %v", attempt, d, p.Max)
		}
	}
}

func TestDelay_NegativeAttemptClamped(t *testing.T) {
	p := backoff.Policy{Base: 200 * time.Millisecond, Max: 10 * time.Second, Factor: 2.0}
	if p.Delay(-3) != p.Delay(0) {
		t.Fatal("negative attempt should equal attempt 0")
	}
}

func TestDelay_JitterIncreasesDelay(t *testing.T) {
	// Run many samples; at least one should differ from the no-jitter value.
	base := backoff.Policy{Base: 500 * time.Millisecond, Max: 30 * time.Second, Factor: 2.0, Jitter: false}
	jittered := backoff.Policy{Base: 500 * time.Millisecond, Max: 30 * time.Second, Factor: 2.0, Jitter: true}

	noJitter := base.Delay(2)
	found := false
	for i := 0; i < 50; i++ {
		if jittered.Delay(2) != noJitter {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("jitter never produced a different delay in 50 samples")
	}
}

func TestDefaultPolicy_ReturnsNonZero(t *testing.T) {
	p := backoff.DefaultPolicy()
	if p.Base <= 0 || p.Max <= 0 || p.Factor <= 1 {
		t.Fatalf("unexpected default policy: %+v", p)
	}
}

func TestSequence_EmitsCorrectCount(t *testing.T) {
	p := backoff.Policy{Base: 10 * time.Millisecond, Max: 1 * time.Second, Factor: 2.0}
	const n = 5
	count := 0
	for range p.Sequence(n) {
		count++
	}
	if count != n {
		t.Fatalf("want %d delays, got %d", n, count)
	}
}

func TestSequence_DelaysAreNonDecreasing(t *testing.T) {
	p := backoff.Policy{Base: 10 * time.Millisecond, Max: 5 * time.Second, Factor: 2.0}
	var prev time.Duration
	for d := range p.Sequence(6) {
		if d < prev {
			t.Fatalf("sequence not non-decreasing: %v < %v", d, prev)
		}
		prev = d
	}
}
