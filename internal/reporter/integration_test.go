package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/reporter"
)

// TestRoundTrip_MultiFlush verifies that multiple flush cycles each emit only
// the records ingested since the previous flush.
func TestRoundTrip_MultiFlush(t *testing.T) {
	var buf bytes.Buffer
	h, _ := history.New("", 100)
	r := reporter.New(&buf, "json", h)

	now := time.Now()

	// cycle 1
	r.Record("a", 80, true, now)
	r.Record("b", 443, true, now)
	if err := r.Flush(); err != nil {
		t.Fatalf("flush 1: %v", err)
	}
	lines1 := countLines(buf.Bytes())
	if lines1 != 2 {
		t.Fatalf("cycle 1: want 2 lines, got %d", lines1)
	}
	buf.Reset()

	// cycle 2 — only one new record
	r.Record("a", 80, false, now.Add(time.Second))
	if err := r.Flush(); err != nil {
		t.Fatalf("flush 2: %v", err)
	}
	lines2 := countLines(buf.Bytes())
	if lines2 != 1 {
		t.Fatalf("cycle 2: want 1 line, got %d", lines2)
	}

	var m map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["host"] != "a" {
		t.Errorf("expected host=a, got %v", m["host"])
	}
}

func countLines(b []byte) int {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return 0
	}
	return len(bytes.Split(b, []byte("\n")))
}
