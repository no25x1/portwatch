// Package aggregator batches scan events over a time window and emits
// a single consolidated report per flush interval.
package aggregator

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Event is a lightweight scan result passed into the aggregator.
type Event struct {
	Host  string
	Port  int
	Open  bool
	Seen  time.Time
}

// Report is the consolidated output emitted after each flush.
type Report struct {
	Window  time.Duration
	Flushed time.Time
	Events  []Event
}

// Aggregator collects events and flushes them on a fixed interval.
type Aggregator struct {
	mu       sync.Mutex
	buf      []Event
	window   time.Duration
	output   chan Report
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates an Aggregator that flushes every window duration.
// The returned channel receives one Report per flush tick.
func New(window time.Duration) (*Aggregator, <-chan Report) {
	if window <= 0 {
		window = 5 * time.Second
	}
	ch := make(chan Report, 8)
	a := &Aggregator{
		window: window,
		output: ch,
		stop:   make(chan struct{}),
	}
	a.wg.Add(1)
	go a.loop()
	return a, ch
}

// Add enqueues an event for the current window.
func (a *Aggregator) Add(e Event) {
	if e.Seen.IsZero() {
		e.Seen = time.Now()
	}
	a.mu.Lock()
	a.buf = append(a.buf, e)
	a.mu.Unlock()
}

// AddFromSnapshot converts a snapshot diff entry into an Event and enqueues it.
func (a *Aggregator) AddFromSnapshot(host string, port int, entry snapshot.Entry) {
	a.Add(Event{
		Host: host,
		Port: port,
		Open: entry.Open,
		Seen: entry.LastSeen,
	})
}

// Stop shuts down the background flush goroutine and closes the output channel.
func (a *Aggregator) Stop() {
	close(a.stop)
	a.wg.Wait()
	close(a.output)
}

func (a *Aggregator) loop() {
	defer a.wg.Done()
	ticker := time.NewTicker(a.window)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.flush()
		case <-a.stop:
			a.flush()
			return
		}
	}
}

func (a *Aggregator) flush() {
	a.mu.Lock()
	events := a.buf
	a.buf = nil
	a.mu.Unlock()
	if len(events) == 0 {
		return
	}
	a.output <- Report{
		Window:  a.window,
		Flushed: time.Now(),
		Events:  events,
	}
}
