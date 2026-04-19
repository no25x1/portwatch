// Package metrics tracks runtime counters for portwatch operations.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	ScansTotal   int64
	AlertsTotal  int64
	ErrorsTotal  int64
	UpCount      int64
	DownCount    int64
	LastScanTime time.Time
}

// Collector accumulates metrics in a thread-safe manner.
type Collector struct {
	mu           sync.Mutex
	scansTotal   int64
	alertsTotal  int64
	errorsTotal  int64
	upCount      int64
	downCount    int64
	lastScanTime time.Time
}

// New returns an initialized Collector.
func New() *Collector {
	return &Collector{}
}

// RecordScan increments the scan counter and records the time.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanTime = time.Now()
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// RecordError increments the error counter.
func (c *Collector) RecordError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorsTotal++
}

// RecordPortUp increments the up counter.
func (c *Collector) RecordPortUp() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.upCount++
}

// RecordPortDown increments the down counter.
func (c *Collector) RecordPortDown() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.downCount++
}

// Snapshot returns a copy of current counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		ScansTotal:   c.scansTotal,
		AlertsTotal:  c.alertsTotal,
		ErrorsTotal:  c.errorsTotal,
		UpCount:      c.upCount,
		DownCount:    c.downCount,
		LastScanTime: c.lastScanTime,
	}
}
