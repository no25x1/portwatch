package poller

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Target represents a host+port pair to monitor.
type Target struct {
	Host string
	Port int
}

// Poller periodically scans targets and fires alerts on state changes.
type Poller struct {
	targets  []Target
	interval time.Duration
	state    *state.Store
	alerter  *alert.Alerter
	opts     scanner.Options
}

// New creates a Poller with the given targets, interval, store, and alerter.
func New(targets []Target, interval time.Duration, s *state.Store, a *alert.Alerter, opts scanner.Options) *Poller {
	return &Poller{
		targets:  targets,
		interval: interval,
		state:    s,
		alerter:  a,
		opts:     opts,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (p *Poller) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.poll()

	for {
		select {
		case <-ticker.C:
			p.poll()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Poller) poll() {
	for _, t := range p.targets {
		open := scanner.CheckPort(t.Host, t.Port, p.opts)
		event, changed := p.state.Update(t.Host, t.Port, open)
		if changed {
			p.alerter.Notify(event)
		}
	}
}
