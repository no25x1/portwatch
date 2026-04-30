package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeEntry(host string, port int, prev, curr string) Entry {
	return Entry{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Host:      host,
		Port:      port,
		Prev:      prev,
		Curr:      curr,
		Source:    "test",
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	l := New(10)
	l.Record(makeEntry("host1", 80, "closed", "open"))
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
}

func TestRecord_Eviction(t *testing.T) {
	l := New(3)
	for i := 0; i < 5; i++ {
		l.Record(makeEntry("h", i, "closed", "open"))
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", l.Len())
	}
	// oldest two should be gone; first remaining port should be 2
	all := l.All()
	if all[0].Port != 2 {
		t.Errorf("expected oldest remaining port=2, got %d", all[0].Port)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New(10)
	l.Record(makeEntry("h", 443, "open", "closed"))
	snap := l.All()
	snap[0].Host = "mutated"
	if l.All()[0].Host == "mutated" {
		t.Error("All() should return an independent copy")
	}
}

func TestFlush_WritesNDJSON(t *testing.T) {
	l := New(10)
	l.Record(makeEntry("web", 8080, "closed", "open"))
	l.Record(makeEntry("db", 5432, "open", "closed"))

	var buf bytes.Buffer
	if err := l.Flush(&buf); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSON lines, got %d", len(lines))
	}
	var e Entry
	if err := json.Unmarshal([]byte(lines[0]), &e); err != nil {
		t.Fatalf("line 0 not valid JSON: %v", err)
	}
	if e.Host != "web" || e.Port != 8080 {
		t.Errorf("unexpected first entry: %+v", e)
	}
}

func TestFlush_ClearsLog(t *testing.T) {
	l := New(10)
	l.Record(makeEntry("h", 22, "closed", "open"))
	var buf bytes.Buffer
	_ = l.Flush(&buf)
	if l.Len() != 0 {
		t.Errorf("expected log to be empty after Flush, got %d", l.Len())
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	l := New(0)
	for i := 0; i < 1001; i++ {
		l.Record(makeEntry("h", i, "closed", "open"))
	}
	if l.Len() != 1000 {
		t.Errorf("expected default cap 1000, got %d", l.Len())
	}
}
