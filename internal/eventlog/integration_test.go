package eventlog_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/eventlog"
)

func TestConcurrentAppend_NeverExceedsMax(t *testing.T) {
	const max = 50
	l := eventlog.New(max)

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Append(eventlog.Entry{
				Timestamp: time.Now(),
				Host:      "10.0.0.1",
				Port:      i % 65535,
				State:     "open",
			})
		}(i)
	}
	wg.Wait()

	if got := l.Len(); got > max {
		t.Errorf("log exceeded max size: got %d, want <= %d", got, max)
	}
}

func TestDrainThenAppend_LogGrowsAgain(t *testing.T) {
	l := eventlog.New(10)
	for i := 0; i < 5; i++ {
		l.Append(eventlog.Entry{Timestamp: time.Now(), Host: "h", Port: i, State: "open"})
	}

	drained := l.Drain()
	if len(drained) != 5 {
		t.Fatalf("expected 5 drained, got %d", len(drained))
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty log after drain")
	}

	l.Append(eventlog.Entry{Timestamp: time.Now(), Host: "h", Port: 9999, State: "closed"})
	if l.Len() != 1 {
		t.Errorf("expected 1 entry after re-append, got %d", l.Len())
	}
}
