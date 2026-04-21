package dedupe_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/dedupe"
)

func TestAllow_FirstEventAlwaysPermitted(t *testing.T) {
	f := dedupe.New(5 * time.Second)
	if !f.Allow("host1", 80, true) {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_SameStateSuppressedWithinCooldown(t *testing.T) {
	f := dedupe.New(5 * time.Second)
	f.Allow("host1", 80, true)
	if f.Allow("host1", 80, true) {
		t.Fatal("expected duplicate event within cooldown to be suppressed")
	}
}

func TestAllow_StateChangeAlwaysPermitted(t *testing.T) {
	f := dedupe.New(5 * time.Second)
	f.Allow("host1", 80, true)
	if !f.Allow("host1", 80, false) {
		t.Fatal("expected state change (open->closed) to be allowed")
	}
}

func TestAllow_AfterCooldownExpires(t *testing.T) {
	var fakeNow time.Time
	f := dedupe.New(100 * time.Millisecond)
	// Inject controllable clock.
	f2 := &struct{ *dedupe.Filter }{dedupe.New(100 * time.Millisecond)}
	_ = f2

	// Use zero cooldown variant to test time-based expiry indirectly.
	f.Allow("host1", 443, true)
	time.Sleep(150 * time.Millisecond)
	_ = fakeNow
	// After cooldown the same state should pass again.
	if !f.Allow("host1", 443, true) {
		t.Fatal("expected event to be allowed after cooldown expires")
	}
}

func TestAllow_ZeroCooldown_SuppressesSameState(t *testing.T) {
	f := dedupe.New(0)
	f.Allow("h", 22, false)
	if f.Allow("h", 22, false) {
		t.Fatal("zero cooldown should still suppress identical consecutive state")
	}
}

func TestAllow_ZeroCooldown_AllowsStateChange(t *testing.T) {
	f := dedupe.New(0)
	f.Allow("h", 22, false)
	if !f.Allow("h", 22, true) {
		t.Fatal("zero cooldown should allow state change")
	}
}

func TestReset_ClearsState(t *testing.T) {
	f := dedupe.New(5 * time.Second)
	f.Allow("host1", 80, true)
	f.Reset()
	if f.Len() != 0 {
		t.Fatalf("expected 0 entries after Reset, got %d", f.Len())
	}
	if !f.Allow("host1", 80, true) {
		t.Fatal("expected event to be allowed after reset")
	}
}

func TestLen_TracksDistinctKeys(t *testing.T) {
	f := dedupe.New(time.Minute)
	f.Allow("a", 80, true)
	f.Allow("b", 80, true)
	f.Allow("a", 443, true)
	if got := f.Len(); got != 3 {
		t.Fatalf("expected 3 distinct keys, got %d", got)
	}
}
