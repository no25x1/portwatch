package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event describes a port state change.
type Event struct {
	Host      string
	Port      int
	PrevOpen  bool
	CurrOpen  bool
	Timestamp time.Time
}

// Notifier sends alerts for port state change events.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify emits an alert for the given event.
func (n *Notifier) Notify(e Event) {
	level := levelFor(e)
	fmt.Fprintf(
		n.out,
		"[%s] %s %s:%d — was_open=%v now_open=%v\n",
		e.Timestamp.Format(time.RFC3339),
		level,
		e.Host,
		e.Port,
		e.PrevOpen,
		e.CurrOpen,
	)
}

// levelFor returns the appropriate alert level for an event.
func levelFor(e Event) Level {
	switch {
	case !e.PrevOpen && e.CurrOpen:
		return LevelInfo // port came up
	case e.PrevOpen && !e.CurrOpen:
		return LevelError // port went down
	default:
		return LevelWarn
	}
}
