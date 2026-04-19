package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/reporter"
)

func newHistory(t *testing.T) *history.History {
	t.Helper()
	h, err := history.New("", 50)
	if err != nil {
		t.Fatalf("history.New: %v", err)
	}
	return h
}

func TestFlush_WritesJSONEntries(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, "json", newHistory(t))

	now := time.Now()
	r.Record("host1", 80, true, now)
	r.Record("host1", 443, false, now)

	if err := r.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, line := range lines {
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			t.Errorf("invalid JSON line %q: %v", line, err)
		}
	}
}

func TestFlush_ResetsBuilder(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, "json", newHistory(t))

	now := time.Now()
	r.Record("host1", 80, true, now)
	_ = r.Flush()
	buf.Reset()

	// second flush with no new records should produce no output
	if err := r.Flush(); err != nil {
		t.Fatalf("second Flush: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output after reset, got %q", buf.String())
	}
}

func TestRecord_PopulatesHistory(t *testing.T) {
	var buf bytes.Buffer
	h := newHistory(t)
	r := reporter.New(&buf, "text", h)

	now := time.Now()
	r.Record("host2", 22, true, now)

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(entries))
	}
	if entries[0].Host != "host2" || entries[0].Port != 22 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}
