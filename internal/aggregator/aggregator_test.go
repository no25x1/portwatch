package aggregator_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/aggregator"
)

func makeEvent(host string, port int, open bool) aggregator.Event {
	return aggregator.Event{
		Host: host,
		Port: port,
		Open: open,
		Seen: time.Now(),
	}
}

func TestAdd_AppearsInReport(t *testing.T) {
	a, ch := aggregator.New(50 * time.Millisecond)
	defer a.Stop()

	a.Add(makeEvent("host-a", 80, true))
	a.Add(makeEvent("host-b", 443, false))

	select {
	case r := <-ch:
		if len(r.Events) != 2 {
			t.Fatalf("expected 2 events, got %d", len(r.Events))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for report")
	}
}

func TestFlush_EmptyWindowProducesNoReport(t *testing.T) {
	a, ch := aggregator.New(40 * time.Millisecond)

	// Let one tick pass without adding anything, then stop.
	time.Sleep(80 * time.Millisecond)
	a.Stop()

	select {
	case r, ok := <-ch:
		if ok {
			t.Fatalf("expected no report for empty window, got %+v", r)
		}
		// channel closed cleanly — expected
	default:
		// nothing buffered — also fine, drain closed channel
		for range ch {
		}
	}
}

func TestStop_FlushesRemainingEvents(t *testing.T) {
	a, ch := aggregator.New(10 * time.Second) // very long window

	a.Add(makeEvent("host-c", 22, true))
	a.Stop() // should flush before returning

	var reports []aggregator.Report
	for r := range ch {
		reports = append(reports, r)
	}
	if len(reports) == 0 {
		t.Fatal("expected at least one report after Stop")
	}
	if len(reports[0].Events) != 1 {
		t.Fatalf("expected 1 event in final report, got %d", len(reports[0].Events))
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	// zero window should be clamped to 5 s internally; just ensure no panic
	a, ch := aggregator.New(0)
	a.Add(makeEvent("host-d", 8080, false))
	a.Stop()
	for range ch {
	}
}

func TestReport_WindowFieldSet(t *testing.T) {
	window := 60 * time.Millisecond
	a, ch := aggregator.New(window)
	defer a.Stop()

	a.Add(makeEvent("host-e", 3000, true))

	select {
	case r := <-ch:
		if r.Window != window {
			t.Fatalf("expected window %v, got %v", window, r.Window)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out")
	}
}
