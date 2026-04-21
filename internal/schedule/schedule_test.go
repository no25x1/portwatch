package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com//portwatch/internal/schedule"
)

func TestNew_DefaultInterval(t *testing.T) {
	s := schedule.New(0)
	if s.Interval() != 30*time.Second {
		t.Fatalf("expected 30s default, got %v", s.Interval())
	}
}

func TestNew_CustomInterval(t *testing.T) {
	s := schedule.New(5 * time.Second)
	if s.Interval() != 5*time.Second {
		t.Fatalf("expected 5s, got %v", s.Interval())
	}
}

func TestRun_FiresImmediately(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s := schedule.New(10 * time.Second)
	ch := s.Run(ctx)

	select {
	case tick, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before first tick")
		}
		if tick.IsZero() {
			t.Fatal("expected non-zero time")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("did not receive immediate tick")
	}
}

func TestRun_TicksAtInterval(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s := schedule.New(200 * time.Millisecond)
	ch := s.Run(ctx)

	var count int
	for range ch {
		count++
		if count >= 3 {
			cancel()
			break
		}
	}
	if count < 3 {
		t.Fatalf("expected at least 3 ticks, got %d", count)
	}
}

func TestRun_ClosesOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := schedule.New(100 * time.Millisecond)
	ch := s.Run(ctx)

	// Drain the immediate tick.
	<-ch
	cancel()

	select {
	case _, ok := <-ch:
		if ok {
			// drain remaining buffered ticks
			for range ch {
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel was not closed after context cancel")
	}
}
