// Package watcher ties together the poller, state tracker, and notifier
// into a single watch loop that runs until the context is cancelled.
package watcher

import (
	"context"
	"log"
	"time"

	"github.com/example/portwatch/internal/notify"
	"github.com/example/portwatch/internal/poller"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning and change notification.
type Watcher struct {
	targets  []poller.Target
	state    *state.Store
	notifier *notify.Dispatcher
	interval time.Duration
	opts     scanner.Options
}

// Config holds watcher configuration.
type Config struct {
	Targets  []poller.Target
	State    *state.Store
	Notifier *notify.Dispatcher
	Interval time.Duration
	Opts     scanner.Options
}

// New creates a Watcher from the provided Config.
func New(cfg Config) *Watcher {
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	return &Watcher{
		targets:  cfg.Targets,
		state:    cfg.State,
		notifier: cfg.Notifier,
		interval: cfg.Interval,
		opts:     cfg.Opts,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.scan(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.scan(ctx)
		}
	}
}

func (w *Watcher) scan(ctx context.Context) {
	for _, t := range w.targets {
		results := scanner.CheckPorts(t.Host, t.Ports, w.opts)
		for _, r := range results {
			ev, changed := w.state.Update(r.Host, r.Port, r.Open)
			if changed {
				if errs := w.notifier.Dispatch(ctx, ev); len(errs) > 0 {
					for _, err := range errs {
						log.Printf("notify error: %v", err)
					}
				}
			}
		}
	}
}
