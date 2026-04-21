// Package pipeline wires together scanning, filtering, deduplication,
// rate-limiting, and notification into a single reusable processing chain.
package pipeline

import (
	"context"
	"log"

	"github.com/yourorg/portwatch/internal/dedupe"
	"github.com/yourorg/portwatch/internal/filter"
	"github.com/yourorg/portwatch/internal/metrics"
	"github.com/yourorg/portwatch/internal/notify"
	"github.com/yourorg/portwatch/internal/ratelimit"
	"github.com/yourorg/portwatch/internal/scanner"
)

// Event represents a port state change emitted by the pipeline.
type Event struct {
	Host   string
	Port   int
	Open   bool
	Prev   bool
}

// Options configures the pipeline behaviour.
type Options struct {
	Filters   []filter.Func
	Dedupe    *dedupe.Deduper
	RateLimit *ratelimit.Limiter
	Metrics   *metrics.Metrics
	Notifier  *notify.Dispatcher
}

// Pipeline processes scan results end-to-end.
type Pipeline struct {
	opts Options
}

// New returns a Pipeline configured with opts.
func New(opts Options) *Pipeline {
	return &Pipeline{opts: opts}
}

// Process takes a slice of scanner.Result values, applies filters, dedup,
// rate-limiting, records metrics, and dispatches notifications.
func (p *Pipeline) Process(ctx context.Context, results []scanner.Result) {
	for _, r := range results {
		ev := Event{Host: r.Host, Port: r.Port, Open: r.Open}

		// Apply filter chain.
		if !p.passesFilters(ev) {
			continue
		}

		// Deduplication check.
		if p.opts.Dedupe != nil && !p.opts.Dedupe.Allow(r.Host, r.Port, r.Open) {
			continue
		}

		// Rate-limit check.
		if p.opts.RateLimit != nil && !p.opts.RateLimit.Allow(r.Host, r.Port) {
			continue
		}

		// Record metrics.
		if p.opts.Metrics != nil {
			p.opts.Metrics.RecordScan(r.Host, r.Port)
			if r.Open {
				p.opts.Metrics.RecordPortUp(r.Host, r.Port)
			} else {
				p.opts.Metrics.RecordPortDown(r.Host, r.Port)
			}
		}

		// Dispatch notifications.
		if p.opts.Notifier != nil {
			if err := p.opts.Notifier.Dispatch(ctx, toNotifyEvent(ev)); err != nil {
				log.Printf("pipeline: notify error: %v", err)
				if p.opts.Metrics != nil {
					p.opts.Metrics.RecordError()
				}
			}
		}
	}
}

func (p *Pipeline) passesFilters(ev Event) bool {
	for _, f := range p.opts.Filters {
		if !f(toFilterEvent(ev)) {
			return false
		}
	}
	return true
}
