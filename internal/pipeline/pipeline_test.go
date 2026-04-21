package pipeline_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/dedupe"
	"github.com/yourorg/portwatch/internal/filter"
	"github.com/yourorg/portwatch/internal/metrics"
	"github.com/yourorg/portwatch/internal/notify"
	"github.com/yourorg/portwatch/internal/pipeline"
	"github.com/yourorg/portwatch/internal/ratelimit"
	"github.com/yourorg/portwatch/internal/scanner"
)

// captureChannel records every event dispatched to it.
type captureChannel struct {
	mu     sync.Mutex
	events []notify.Event
}

func (c *captureChannel) Send(_ context.Context, e notify.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, e)
	return nil
}

func (c *captureChannel) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.events)
}

func makeResults(host string, port int, open bool) []scanner.Result {
	return []scanner.Result{{Host: host, Port: port, Open: open}}
}

func TestProcess_DispatchesEvent(t *testing.T) {
	ch := &captureChannel{}
	p := pipeline.New(pipeline.Options{
		Notifier: notify.New(ch),
	})
	p.Process(context.Background(), makeResults("localhost", 8080, true))
	if ch.count() != 1 {
		t.Fatalf("expected 1 event, got %d", ch.count())
	}
}

func TestProcess_FilterDropsEvent(t *testing.T) {
	ch := &captureChannel{}
	p := pipeline.New(pipeline.Options{
		Filters:  []filter.Func{filter.OnlyClosed},
		Notifier: notify.New(ch),
	})
	// Open event should be dropped by OnlyClosed filter.
	p.Process(context.Background(), makeResults("localhost", 8080, true))
	if ch.count() != 0 {
		t.Fatalf("expected 0 events, got %d", ch.count())
	}
}

func TestProcess_DedupeSuppress(t *testing.T) {
	ch := &captureChannel{}
	p := pipeline.New(pipeline.Options{
		Dedupe:   dedupe.New(10 * time.Second),
		Notifier: notify.New(ch),
	})
	res := makeResults("host1", 9000, true)
	p.Process(context.Background(), res)
	p.Process(context.Background(), res) // same state — should be suppressed
	if ch.count() != 1 {
		t.Fatalf("expected 1 event after dedup, got %d", ch.count())
	}
}

func TestProcess_MetricsRecorded(t *testing.T) {
	m := metrics.New()
	p := pipeline.New(pipeline.Options{
		Metrics: m,
	})
	p.Process(context.Background(), makeResults("host1", 443, true))
	snap := m.Snapshot()
	if snap.TotalScans == 0 {
		t.Fatal("expected TotalScans > 0")
	}
}

func TestProcess_RateLimitSuppressesSecondAlert(t *testing.T) {
	ch := &captureChannel{}
	p := pipeline.New(pipeline.Options{
		RateLimit: ratelimit.New(10 * time.Second),
		Notifier:  notify.New(ch),
	})
	res := makeResults("host2", 22, false)
	p.Process(context.Background(), res)
	p.Process(context.Background(), res)
	if ch.count() != 1 {
		t.Fatalf("expected 1 event due to rate-limit, got %d", ch.count())
	}
}
