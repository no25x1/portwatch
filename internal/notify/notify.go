// Package notify dispatches alerts to one or more notification channels.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event describes a port state change to be dispatched.
type Event struct {
	Host      string
	Port      int
	PrevState string
	CurrState string
	OccurredAt time.Time
}

// Channel is anything that can receive an Event.
type Channel interface {
	Send(e Event) error
}

// Dispatcher fans out events to all registered channels.
type Dispatcher struct {
	channels []Channel
}

// New returns a Dispatcher with the given channels.
func New(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Dispatch sends e to every registered channel, collecting errors.
func (d *Dispatcher) Dispatch(e Event) []error {
	var errs []error
	for _, ch := range d.channels {
		if err := ch.Send(e); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// StdoutChannel writes a human-readable line to an io.Writer (default os.Stdout).
type StdoutChannel struct {
	w io.Writer
}

// NewStdoutChannel creates a StdoutChannel writing to w. If w is nil, os.Stdout is used.
func NewStdoutChannel(w io.Writer) *StdoutChannel {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutChannel{w: w}
}

// Send writes the event as a formatted line.
func (s *StdoutChannel) Send(e Event) error {
	_, err := fmt.Fprintf(s.w, "[%s] %s:%d  %s -> %s\n",
		e.OccurredAt.Format(time.RFC3339), e.Host, e.Port, e.PrevState, e.CurrState)
	return err
}
