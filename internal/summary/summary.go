// Package summary provides aggregation of port scan results
// into a structured report for display or export.
package summary

import (
	"fmt"
	"time"
)

// PortStatus represents the last known status of a single port.
type PortStatus struct {
	Host    string
	Port    int
	Open    bool
	LastSeen time.Time
}

// Report holds an aggregated snapshot of all monitored ports.
type Report struct {
	GeneratedAt time.Time
	Entries     []PortStatus
}

// TotalOpen returns the number of open ports in the report.
func (r *Report) TotalOpen() int {
	count := 0
	for _, e := range r.Entries {
		if e.Open {
			count++
		}
	}
	return count
}

// TotalClosed returns the number of closed ports in the report.
func (r *Report) TotalClosed() int {
	return len(r.Entries) - r.TotalOpen()
}

// Key returns a unique string key for a PortStatus entry.
func Key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Builder accumulates PortStatus entries and builds a Report.
type Builder struct {
	entries map[string]PortStatus
}

// New creates a new Builder.
func New() *Builder {
	return &Builder{entries: make(map[string]PortStatus)}
}

// Upsert adds or updates a port status entry.
func (b *Builder) Upsert(host string, port int, open bool, lastSeen time.Time) {
	k := Key(host, port)
	b.entries[k] = PortStatus{
		Host:     host,
		Port:     port,
		Open:     open,
		LastSeen: lastSeen,
	}
}

// Build returns a Report from the accumulated entries.
func (b *Builder) Build() Report {
	r := Report{GeneratedAt: time.Now()}
	for _, e := range b.entries {
		r.Entries = append(r.Entries, e)
	}
	return r
}
