// Package audit provides a structured audit trail for port state transitions,
// recording who observed a change, when, and what the previous and current
// states were. Entries are kept in memory up to a configurable maximum and
// can be flushed to an io.Writer as newline-delimited JSON.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// Entry captures a single auditable state-change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Prev      string    `json:"prev_state"`
	Curr      string    `json:"curr_state"`
	Source    string    `json:"source"`
}

// Log is a bounded, thread-safe audit log.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New returns a Log that retains at most maxSize entries.
// If maxSize <= 0 it defaults to 1000.
func New(maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &Log{maxSize: maxSize}
}

// Record appends an Entry to the log, evicting the oldest entry when full.
func (l *Log) Record(e Entry) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, e)
}

// All returns a snapshot of all current entries in chronological order.
func (l *Log) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Flush writes every entry as newline-delimited JSON to w, then clears the log.
// Returns the first encoding error encountered, if any.
func (l *Log) Flush(w io.Writer) error {
	l.mu.Lock()
	entries := l.entries
	l.entries = nil
	l.mu.Unlock()

	enc := json.NewEncoder(w)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("audit: encode entry: %w", err)
		}
	}
	return nil
}

// Len returns the current number of stored entries.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
