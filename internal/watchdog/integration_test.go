package watchdog_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/watchdog"
)

// slowRunner simulates a runner that takes a moment to start.
type slowRunner struct {
	fakeRunner
	startDelay time.Duration
}

func (s *slowRunner) Start(ctx context.Context) error {
	select {
	case <-time.After(s.startDelay):
	case <-ctx.Done():
	}
	s.starts.Add(1)
	return nil
}

func TestWatchdog_SlowRunnerRestart(t *testing.T) {
	m := metrics.New()
	for i := 0; i < 5; i++ {
		m.RecordError()
	}

	runner := &slowRunner{startDelay: 5 * time.Millisecond}
	wd := watchdog.New(runner, m, watchdog.Options{
		CheckInterval: 20 * time.Millisecond,
		MaxErrors:     5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	wd.Run(ctx)

	if wd.Restarts() < 1 {
		t.Fatalf("expected restart, got %d", wd.Restarts())
	}
}

func TestWatchdog_CancelStopsLoop(t *testing.T) {
	m := metrics.New()
	runner := &fakeRunner{}
	wd := watchdog.New(runner, m, watchdog.Options{
		CheckInterval: 50 * time.Millisecond,
		MaxErrors:     1,
	})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		wd.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watchdog did not stop after context cancellation")
	}
}
