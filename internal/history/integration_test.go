package history_test

import (
	"path/filepath"
	"testing"

	"github.com/yourorg/portwatch/internal/history"
)

// TestRoundTrip verifies that entries survive a full save/load cycle.
func TestRoundTrip(t *testing.T) {
	h := history.New(100)
	events := []struct {
		host string
		port int
		open bool
	}{
		{"alpha", 80, true},
		{"beta", 443, false},
		{"gamma", 22, true},
	}
	for _, e := range events {
		h.Record(e.host, e.port, e.open)
	}

	path := filepath.Join(t.TempDir(), "rt.json")
	if err := h.SaveJSON(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	h2 := history.New(100)
	if err := h2.LoadJSON(path); err != nil {
		t.Fatalf("load: %v", err)
	}

	all := h2.All()
	if len(all) != len(events) {
		t.Fatalf("expected %d entries, got %d", len(events), len(all))
	}
	for i, e := range events {
		if all[i].Host != e.host || all[i].Port != e.port || all[i].Open != e.open {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, all[i], e)
		}
	}
}

// TestNew_DefaultMaxSize ensures zero/negative maxSize falls back to default.
func TestNew_DefaultMaxSize(t *testing.T) {
	h := history.New(0)
	for i := 0; i < 600; i++ {
		h.Record("host", i, true)
	}
	all := h.All()
	if len(all) != 500 {
		t.Errorf("expected default max 500, got %d", len(all))
	}
}
