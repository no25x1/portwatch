package ratelimit

import (
	"testing"
	"time"
)

func newFakeLimiter(cooldown time.Duration) (*Limiter, *time.Time) {
	ts := time.Now()
	l := New(cooldown)
	l.now = func() time.Time { return ts }
	return l, &ts
}

func TestAllow_FirstCallAlwaysPermitted(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if !l.Allow("host1:80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldown_Suppressed(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("host1:80")
	if l.Allow("host1:80") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldownExpires_Permitted(t *testing.T) {
	l, ts := newFakeLimiter(5 * time.Second)
	l.Allow("host1:80")
	*ts = ts.Add(6 * time.Second)
	if !l.Allow("host1:80") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_ZeroCooldown_AlwaysPermitted(t *testing.T) {
	l := New(0)
	for i := 0; i < 5; i++ {
		if !l.Allow("host1:443") {
			t.Fatalf("expected call %d to be allowed with zero cooldown", i)
		}
	}
}

func TestAllow_DifferentKeys_IndependentWindows(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("host1:80")
	if !l.Allow("host2:80") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("host1:80")
	l.Reset("host1:80")
	if !l.Allow("host1:80") {
		t.Fatal("expected call after Reset to be allowed")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if l.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", l.Len())
	}
	l.Allow("host1:80")
	l.Allow("host2:443")
	if l.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", l.Len())
	}
	l.Reset("host1:80")
	if l.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", l.Len())
	}
}
