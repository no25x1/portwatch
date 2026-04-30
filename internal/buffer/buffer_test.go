package buffer

import (
	"testing"
	"time"
)

func makeEntry(host string, port int, open bool) Entry {
	return Entry{Host: host, Port: port, Open: open, Timestamp: time.Now()}
}

func TestPush_AppendsEntry(t *testing.T) {
	b := New(4)
	b.Push(makeEntry("host1", 80, true))
	if b.Len() != 1 {
		t.Fatalf("expected len 1, got %d", b.Len())
	}
}

func TestAll_ReturnsInsertionOrder(t *testing.T) {
	b := New(4)
	b.Push(makeEntry("a", 80, true))
	b.Push(makeEntry("b", 443, false))
	b.Push(makeEntry("c", 8080, true))

	all := b.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	if all[0].Host != "a" || all[1].Host != "b" || all[2].Host != "c" {
		t.Errorf("unexpected order: %v", all)
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	b := New(3)
	b.Push(makeEntry("first", 1, true))
	b.Push(makeEntry("second", 2, true))
	b.Push(makeEntry("third", 3, true))
	b.Push(makeEntry("fourth", 4, true)) // evicts "first"

	all := b.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", len(all))
	}
	if all[0].Host != "second" {
		t.Errorf("expected oldest to be 'second', got %q", all[0].Host)
	}
	if all[2].Host != "fourth" {
		t.Errorf("expected newest to be 'fourth', got %q", all[2].Host)
	}
}

func TestDrain_ClearsBuffer(t *testing.T) {
	b := New(4)
	b.Push(makeEntry("x", 22, true))
	b.Push(makeEntry("y", 80, false))

	drained := b.Drain()
	if len(drained) != 2 {
		t.Fatalf("expected 2 drained entries, got %d", len(drained))
	}
	if b.Len() != 0 {
		t.Errorf("expected buffer empty after drain, got len %d", b.Len())
	}
}

func TestNew_PanicOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero capacity")
		}
	}()
	New(0)
}

func TestAll_IsSnapshot(t *testing.T) {
	b := New(4)
	b.Push(makeEntry("snap", 9090, true))

	snap := b.All()
	b.Push(makeEntry("after", 9091, false))

	if len(snap) != 1 {
		t.Errorf("snapshot should not reflect later pushes, got len %d", len(snap))
	}
}
