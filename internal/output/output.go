// Package output formats and writes port event reports.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Format controls how events are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Event represents a single port state change to be reported.
type Event struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	State     string    `json:"state"`
	PrevState string    `json:"prev_state"`
	Timestamp time.Time `json:"timestamp"`
}

// Writer writes formatted events to an io.Writer.
type Writer struct {
	format Format
	out    io.Writer
}

// New creates a Writer with the given format. Defaults to stdout.
func New(format Format, out io.Writer) *Writer {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Writer{format: format, out: out}
}

// Write formats and emits a single event.
func (w *Writer) Write(e Event) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(e)
	default:
		return w.writeText(e)
	}
}

func (w *Writer) writeText(e Event) error {
	_, err := fmt.Fprintf(w.out, "[%s] %s:%d  %s -> %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Host, e.Port,
		e.PrevState, e.State,
	)
	return err
}

func (w *Writer) writeJSON(e Event) error {
	enc := json.NewEncoder(w.out)
	return enc.Encode(e)
}
