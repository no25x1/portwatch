package debounce_test

import (
	"testing"
	"time"

	"portwatch/internal/debounce"
)

const quiet = 50 * time.Millisecond

func TestSubmit_ForwardsAfterQuiet(t *testing.T) {
	d, ch := debounce.New(quiet)

	d.Submit(debounce.Event{Host: "host1", Port: 80, Open: true})

	select {
	case ev := <-ch:
		if ev.Host != "host1" || ev.Port != 80 || !ev.Open {
			t.Fatalf("unexpected event: %+v", ev)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for forwarded event")
	}
}

func TestSubmit_ResetsTimerOnRepeat(t *testing.T) {
	d, ch := debounce.New(quiet)

	// Submit the same port twice quickly; only one event should arrive.
	d.Submit(debounce.Event{Host: "host1", Port: 443, Open: false})
	time.Sleep(quiet / 2)
	d.Submit(debounce.Event{Host: "host1", Port: 443, Open: true})

	select {
	case ev := <-ch:
		// The second submission wins.
		if !ev.Open {
			t.Fatalf("expected final Open=true, got false")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for forwarded event")
	}

	// Ensure no second event arrives.
	select {
	case extra := <-ch:
		t.Fatalf("unexpected extra event: %+v", extra)
	case <-time.After(quiet * 2):
		// good – only one event forwarded
	}
}

func TestSubmit_DifferentPortsAreIndependent(t *testing.T) {
	d, ch := debounce.New(quiet)

	d.Submit(debounce.Event{Host: "h", Port: 80, Open: true})
	d.Submit(debounce.Event{Host: "h", Port: 443, Open: true})

	received := map[int]bool{}
	deadline := time.After(300 * time.Millisecond)
	for len(received) < 2 {
		select {
		case ev := <-ch:
			received[ev.Port] = true
		case <-deadline:
			t.Fatalf("only received events for ports: %v", received)
		}
	}
}

func TestPending_CountsInFlight(t *testing.T) {
	d, _ := debounce.New(10 * time.Second) // long window so nothing fires

	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending, got %d", d.Pending())
	}

	d.Submit(debounce.Event{Host: "h", Port: 22, Open: false})
	d.Submit(debounce.Event{Host: "h", Port: 8080, Open: true})

	if got := d.Pending(); got != 2 {
		t.Fatalf("expected 2 pending, got %d", got)
	}
}
