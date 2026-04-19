package state

import (
	"os"
	"testing"
	"time"
)

func TestUpdate_DetectsChange(t *testing.T) {
	s, _ := New("")
	key := PortKey{Host: "localhost", Port: 8080}
	now := time.Now()

	// First update — no prior state, always changed.
	if !s.Update(key, true, now) {
		t.Fatal("expected changed=true on first update")
	}

	// Same state — no change.
	if s.Update(key, true, now.Add(time.Second)) {
		t.Fatal("expected changed=false when state unchanged")
	}

	// State flipped — changed.
	if !s.Update(key, false, now.Add(2*time.Second)) {
		t.Fatal("expected changed=true when state flipped")
	}
}

func TestUpdate_LastSeen(t *testing.T) {
	s, _ := New("")
	key := PortKey{Host: "localhost", Port: 9090}
	t0 := time.Now()

	s.Update(key, true, t0)
	st, _ := s.Get(key)
	if !st.LastSeen.Equal(t0) {
		t.Fatalf("expected LastSeen=%v got %v", t0, st.LastSeen)
	}

	// Port goes down — LastSeen should remain t0.
	s.Update(key, false, t0.Add(time.Minute))
	st, _ = s.Get(key)
	if !st.LastSeen.Equal(t0) {
		t.Fatalf("LastSeen should not update when port is down, got %v", st.LastSeen)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s, _ := New("")
	_, ok := s.Get(PortKey{Host: "x", Port: 1})
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	s1, _ := New(tmp.Name())
	key := PortKey{Host: "192.168.1.1", Port: 443}
	now := time.Now().Truncate(time.Second)
	s1.Update(key, true, now)
	if err := s1.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2, err := New(tmp.Name())
	if err != nil {
		t.Fatalf("New (load): %v", err)
	}
	st, ok := s2.Get(key)
	if !ok {
		t.Fatal("expected key to be loaded")
	}
	if !st.Open {
		t.Error("expected Open=true after reload")
	}
}
