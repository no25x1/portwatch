package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/watchdog"
)

// fakeRunner counts Start/Stop calls.
type fakeRunner struct {
	starts atomic.Int32
	stops  atomic.Int32
}

func (f *fakeRunner) Start(_ context.Context) error {
	f.starts.Add(1)
	return nil
}

func (f *fakeRunner) Stop() {
	f.stops.Add(1)
}

func TestWatchdog_NoRestartBelowThreshold(t *testing.T) {
	m := metrics.New()
	// Record fewer errors than the threshold.
	m.RecordError()
	m.RecordError()

	runner := &fakeRunner{}
	wd := watchdog.New(runner, m, watchdog.Options{
		CheckInterval: 20 * time.Millisecond,
		MaxErrors:     5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	wd.Run(ctx)

	if wd.Restarts() != 0 {
		t.Fatalf("expected 0 restarts, got %d", wd.Restarts())
	}
}

func TestWatchdog_RestartsOnThreshold(t *testing.T) {
	m := metrics.New()
	for i := 0; i < 5; i++ {
		m.RecordError()
	}

	runner := &fakeRunner{}
	wd := watchdog.New(runner, m, watchdog.Options{
		CheckInterval: 20 * time.Millisecond,
		MaxErrors:     5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	wd.Run(ctx)

	if wd.Restarts() < 1 {
		t.Fatalf("expected at least 1 restart, got %d", wd.Restarts())
	}
	if runner.stops.Load() < 1 {
		t.Fatal("expected runner.Stop to be called at least once")
	}
}

func TestWatchdog_RestartsIncrement(t *testing.T) {
	m := metrics.New()
	for i := 0; i < 10; i++ {
		m.RecordError()
	}

	runner := &fakeRunner{}
	wd := watchdog.New(runner, m, watchdog.Options{
		CheckInterval: 15 * time.Millisecond,
		MaxErrors:     5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	wd.Run(ctx)

	if wd.Restarts() == 0 {
		t.Fatal("expected at least one restart")
	}
}
