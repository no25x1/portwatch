// Package watchdog provides a self-healing mechanism that monitors the
// poller health and restarts it when the circuit breaker opens or scan
// errors exceed a configurable threshold.
package watchdog

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Runner is anything that can be started and stopped.
type Runner interface {
	Start(ctx context.Context) error
	Stop()
}

// Options configures the watchdog behaviour.
type Options struct {
	// CheckInterval is how often the watchdog evaluates runner health.
	CheckInterval time.Duration
	// MaxErrors is the number of consecutive scan errors before a restart.
	MaxErrors int64
	// Logger receives diagnostic messages; defaults to log.Default().
	Logger *log.Logger
}

func defaults(o Options) Options {
	if o.CheckInterval <= 0 {
		o.CheckInterval = 15 * time.Second
	}
	if o.MaxErrors <= 0 {
		o.MaxErrors = 5
	}
	if o.Logger == nil {
		o.Logger = log.Default()
	}
	return o
}

// Watchdog supervises a Runner and restarts it when unhealthy.
type Watchdog struct {
	opts    Options
	runner  Runner
	metrics *metrics.Metrics

	mu       sync.Mutex
	restarts int
}

// New creates a Watchdog that supervises runner using m for health signals.
func New(runner Runner, m *metrics.Metrics, opts Options) *Watchdog {
	return &Watchdog{
		opts:   defaults(opts),
		runner: runner,
		metrics: m,
	}
}

// Run starts the supervision loop and blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.opts.CheckInterval)
	defer ticker.Stop()

	var lastErrors int64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snap := w.metrics.Snapshot()
			if snap.Errors-lastErrors >= w.opts.MaxErrors {
				w.opts.Logger.Printf("watchdog: error threshold reached (%d), restarting runner", snap.Errors)
				w.restart(ctx)
				lastErrors = snap.Errors
			}
		}
	}
}

func (w *Watchdog) restart(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.runner.Stop()
	w.restarts++

	go func() {
		if err := w.runner.Start(ctx); err != nil && ctx.Err() == nil {
			w.opts.Logger.Printf("watchdog: runner exited with error: %v", err)
		}
	}()
}

// Restarts returns the total number of restarts performed.
func (w *Watchdog) Restarts() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.restarts
}
