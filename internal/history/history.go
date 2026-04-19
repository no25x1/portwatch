// Package history records port event history for reporting and diagnostics.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port event.
type Entry struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Open      bool      `json:"open"`
	Timestamp time.Time `json:"timestamp"`
}

// History stores a bounded list of port events.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a History with the given maximum number of entries.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &History{maxSize: maxSize}
}

// Record appends a new entry, evicting the oldest if at capacity.
func (h *History) Record(host string, port int, open bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Entry{Host: host, Port: port, Open: open, Timestamp: time.Now().UTC()}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

// All returns a copy of all recorded entries.
func (h *History) All() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// SaveJSON writes all entries as JSON to the given file path.
func (h *History) SaveJSON(path string) error {
	h.mu.Lock()
	data, err := json.MarshalIndent(h.entries, "", "  ")
	h.mu.Unlock()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadJSON reads entries from a JSON file, replacing current entries.
func (h *History) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	h.mu.Lock()
	h.entries = entries
	h.mu.Unlock()
	return nil
}
