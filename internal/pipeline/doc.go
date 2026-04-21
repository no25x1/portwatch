// Package pipeline provides an end-to-end processing chain for portwatch scan
// results.
//
// A Pipeline applies, in order:
//
//  1. A configurable filter chain (see internal/filter).
//  2. Deduplication to suppress repeated events for the same port state
//     (see internal/dedupe).
//  3. Per-host/port rate-limiting to avoid alert storms
//     (see internal/ratelimit).
//  4. Metrics recording (see internal/metrics).
//  5. Notification dispatch to one or more output channels
//     (see internal/notify).
//
// Typical usage:
//
//	p := pipeline.New(pipeline.Options{
//		Filters:   []filter.Func{filter.OnlyOpen},
//		Dedupe:    dedupe.New(30 * time.Second),
//		RateLimit: ratelimit.New(5 * time.Minute),
//		Metrics:   metrics.New(),
//		Notifier:  notify.New(channels...),
//	})
//	p.Process(ctx, results)
package pipeline
