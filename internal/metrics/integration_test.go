package metrics_test

import (
	"testing"

	"github.com/user/portwatch/internal/metrics"
)

// TestFullWorkflowCounters simulates a polling cycle recording mixed events.
func TestFullWorkflowCounters(t *testing.T) {
	c := metrics.New()

	// Simulate two scan cycles.
	for i := 0; i < 2; i++ {
		c.RecordScan()
		c.RecordPortUp()
		c.RecordPortUp()
		c.RecordPortDown()
		c.RecordAlert()
		c.RecordAlert()
		c.RecordError()
	}

	snap := c.Snapshot()

	if snap.ScansTotal != 2 {
		t.Errorf("ScansTotal: want 2, got %d", snap.ScansTotal)
	}
	if snap.UpCount != 4 {
		t.Errorf("UpCount: want 4, got %d", snap.UpCount)
	}
	if snap.DownCount != 2 {
		t.Errorf("DownCount: want 2, got %d", snap.DownCount)
	}
	if snap.AlertsTotal != 4 {
		t.Errorf("AlertsTotal: want 4, got %d", snap.AlertsTotal)
	}
	if snap.ErrorsTotal != 2 {
		t.Errorf("ErrorsTotal: want 2, got %d", snap.ErrorsTotal)
	}
	if snap.LastScanTime.IsZero() {
		t.Error("LastScanTime should not be zero")
	}
}
