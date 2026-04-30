// Package debounce delays forwarding of repeated events until a quiet period
// has elapsed, reducing alert noise during flapping port states.
package debounce

import (
	"sync"
	"time"
)

// Event represents a minimal port-state event used by the debouncer.
type Event struct {
	Host  string
	Port  int
	Open  bool
}

// key uniquely identifies a host+port pair.
type key struct {
	host string
	port int
}

type pending struct {
	event Event
	timer *time.Timer
}

// Debouncer holds back events until the port state has been stable for
// at least Quiet duration, then forwards them to the output channel.
type Debouncer struct {
	quiet   time.Duration
	out     chan Event
	mu      sync.Mutex
	pending map[key]*pending
}

// New creates a Debouncer with the given quiet window.
// Events are forwarded on the returned read-only channel.
func New(quiet time.Duration) (*Debouncer, <-chan Event) {
	ch := make(chan Event, 64)
	d := &Debouncer{
		quiet:   quiet,
		out:     ch,
		pending: make(map[key]*pending),
	}
	return d, ch
}

// Submit accepts an event. If a pending timer already exists for the same
// host+port it is reset; otherwise a new timer is started.
func (d *Debouncer) Submit(e Event) {
	k := key{host: e.Host, port: e.Port}

	d.mu.Lock()
	defer d.mu.Unlock()

	if p, ok := d.pending[k]; ok {
		p.event = e
		p.timer.Reset(d.quiet)
		return
	}

	p := &pending{event: e}
	p.timer = time.AfterFunc(d.quiet, func() {
		d.mu.Lock()
		ev := d.pending[k].event
		delete(d.pending, k)
		d.mu.Unlock()
		d.out <- ev
	})
	d.pending[k] = p
}

// Pending returns the number of events currently waiting to be forwarded.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}
