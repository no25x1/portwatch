// Package snapshot captures and compares point-in-time port state
// across all monitored hosts, enabling diff-based change detection.
package snapshot

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Entry holds the observed state of a single host:port at a moment in time.
type Entry struct {
	Host      string
	Port      int
	Open      bool
	Timestamp time.Time
}

// Key returns a canonical string key for an Entry.
func Key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Snapshot is an immutable point-in-time view of port states.
type Snapshot struct {
	capturedAt time.Time
	entries    map[string]Entry
}

// CapturedAt returns when the snapshot was taken.
func (s *Snapshot) CapturedAt() time.Time { return s.capturedAt }

// Get returns the Entry for the given host:port, and whether it was found.
func (s *Snapshot) Get(host string, port int) (Entry, bool) {
	e, ok := s.entries[Key(host, port)]
	return e, ok
}

// All returns all entries sorted by key for deterministic iteration.
func (s *Snapshot) All() []Entry {
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return Key(out[i].Host, out[i].Port) < Key(out[j].Host, out[j].Port)
	})
	return out
}

// Diff returns entries whose Open state differs between s and next.
func (s *Snapshot) Diff(next *Snapshot) []Entry {
	var changed []Entry
	for k, ne := range next.entries {
		if oe, ok := s.entries[k]; !ok || oe.Open != ne.Open {
			changed = append(changed, ne)
		}
	}
	sort.Slice(changed, func(i, j int) bool {
		return Key(changed[i].Host, changed[i].Port) < Key(changed[j].Host, changed[j].Port)
	})
	return changed
}

// Builder accumulates scan results and produces a Snapshot.
type Builder struct {
	mu      sync.Mutex
	entries map[string]Entry
}

// NewBuilder returns a ready-to-use Builder.
func NewBuilder() *Builder {
	return &Builder{entries: make(map[string]Entry)}
}

// Record adds or updates an entry in the builder.
func (b *Builder) Record(host string, port int, open bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[Key(host, port)] = Entry{
		Host:      host,
		Port:      port,
		Open:      open,
		Timestamp: time.Now(),
	}
}

// Build finalises the builder and returns an immutable Snapshot.
func (b *Builder) Build() *Snapshot {
	b.mu.Lock()
	defer b.mu.Unlock()
	cp := make(map[string]Entry, len(b.entries))
	for k, v := range b.entries {
		cp[k] = v
	}
	return &Snapshot{capturedAt: time.Now(), entries: cp}
}
