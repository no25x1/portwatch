// Package eventlog provides a bounded, thread-safe append-only log of
// port-state-change events that can be queried and drained by consumers.
package eventlog

import (
	"sync"
	"time"
)

// Entry holds a single logged event.
type Entry struct {
	Timestamp time.Time
	Host      string
	Port      int
	State     string // "open" | "closed"
}

// Log is a bounded, thread-safe event log.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New returns a Log that retains at most maxSize entries.
// If maxSize is <= 0 it defaults to 1000.
func New(maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &Log{maxSize: maxSize}
}

// Append adds an entry to the log, evicting the oldest entry when the log is full.
func (l *Log) Append(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, e)
}

// All returns a shallow copy of all current entries in insertion order.
func (l *Log) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Drain removes and returns all entries, leaving the log empty.
func (l *Log) Drain() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := l.entries
	l.entries = nil
	return out
}

// Len returns the current number of entries.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
