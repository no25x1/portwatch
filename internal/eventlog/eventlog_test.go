package eventlog

import (
	"testing"
	"time"
)

func makeEntry(host string, port int, state string) Entry {
	return Entry{Timestamp: time.Now(), Host: host, Port: port, State: state}
}

func TestAppend_StoresEntry(t *testing.T) {
	l := New(10)
	l.Append(makeEntry("host-a", 80, "open"))
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New(10)
	l.Append(makeEntry("host-a", 80, "open"))
	l.Append(makeEntry("host-b", 443, "closed"))

	entries := l.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// mutating the returned slice must not affect the log
	entries[0].Host = "mutated"
	if l.All()[0].Host != "host-a" {
		t.Error("All() returned a reference, not a copy")
	}
}

func TestAppend_Eviction(t *testing.T) {
	l := New(3)
	for i := 0; i < 5; i++ {
		l.Append(makeEntry("h", i, "open"))
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", l.Len())
	}
	// oldest entries (ports 0,1) should be gone; remaining are ports 2,3,4
	all := l.All()
	if all[0].Port != 2 {
		t.Errorf("expected oldest remaining port=2, got %d", all[0].Port)
	}
}

func TestDrain_ClearsLog(t *testing.T) {
	l := New(10)
	l.Append(makeEntry("host-a", 22, "open"))
	l.Append(makeEntry("host-a", 22, "closed"))

	drained := l.Drain()
	if len(drained) != 2 {
		t.Fatalf("expected 2 drained entries, got %d", len(drained))
	}
	if l.Len() != 0 {
		t.Errorf("expected log to be empty after Drain, got %d", l.Len())
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	l := New(0)
	for i := 0; i < 1001; i++ {
		l.Append(makeEntry("h", i, "open"))
	}
	if l.Len() != 1000 {
		t.Errorf("expected default max 1000, got %d", l.Len())
	}
}
