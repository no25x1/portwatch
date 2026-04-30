package audit_test

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"

	"github.com/yourorg/portwatch/internal/audit"
)

func TestConcurrentRecord_NeverExceedsMax(t *testing.T) {
	const max = 50
	l := audit.New(max)

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			l.Record(audit.Entry{
				Host:   "h",
				Port:   port,
				Prev:   "closed",
				Curr:   "open",
				Source: "test",
			})
		}(i)
	}
	wg.Wait()

	if n := l.Len(); n > max {
		t.Errorf("log grew beyond max: got %d, want <= %d", n, max)
	}
}

func TestFlushThenRecord_LogGrowsAgain(t *testing.T) {
	l := audit.New(10)
	for i := 0; i < 5; i++ {
		l.Record(audit.Entry{Host: "h", Port: i, Prev: "closed", Curr: "open", Source: "t"})
	}
	var buf bytes.Buffer
	if err := l.Flush(&buf); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty log after flush")
	}
	l.Record(audit.Entry{Host: "h2", Port: 9000, Prev: "open", Curr: "closed", Source: "t"})
	if l.Len() != 1 {
		t.Errorf("expected 1 entry after re-record, got %d", l.Len())
	}

	// verify the flushed bytes are valid NDJSON
	dec := json.NewDecoder(&buf)
	var count int
	for dec.More() {
		var e audit.Entry
		if err := dec.Decode(&e); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		count++
	}
	if count != 5 {
		t.Errorf("expected 5 flushed entries, got %d", count)
	}
}
