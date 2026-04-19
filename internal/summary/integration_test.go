package summary_test

import (
	"testing"
	"time"

	"portwatch/internal/summary"
)

// TestRoundTrip verifies that multiple upserts across hosts produce
// a coherent report with correct open/closed counts.
func TestRoundTrip_MultiHost(t *testing.T) {
	b := summary.New()
	now := time.Now()

	hosts := []struct {
		host string
		port int
		open bool
	}{
		{"host-a", 22, true},
		{"host-a", 80, true},
		{"host-a", 443, false},
		{"host-b", 22, false},
		{"host-b", 3306, true},
	}

	for _, h := range hosts {
		b.Upsert(h.host, h.port, h.open, now)
	}

	r := b.Build()

	if len(r.Entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(r.Entries))
	}
	if r.TotalOpen() != 3 {
		t.Errorf("expected 3 open ports, got %d", r.TotalOpen())
	}
	if r.TotalClosed() != 2 {
		t.Errorf("expected 2 closed ports, got %d", r.TotalClosed())
	}
	if r.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}
