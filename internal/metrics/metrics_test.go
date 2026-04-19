package metrics_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecordScan_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordScan()
	c.RecordScan()
	snap := c.Snapshot()
	if snap.ScansTotal != 2 {
		t.Fatalf("expected 2 scans, got %d", snap.ScansTotal)
	}
	if snap.LastScanTime.IsZero() {
		t.Fatal("expected LastScanTime to be set")
	}
}

func TestRecordAlert_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordAlert()
	if c.Snapshot().AlertsTotal != 1 {
		t.Fatal("expected 1 alert")
	}
}

func TestRecordError_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordError()
	c.RecordError()
	if c.Snapshot().ErrorsTotal != 2 {
		t.Fatal("expected 2 errors")
	}
}

func TestRecordPortUpDown(t *testing.T) {
	c := metrics.New()
	c.RecordPortUp()
	c.RecordPortUp()
	c.RecordPortDown()
	snap := c.Snapshot()
	if snap.UpCount != 2 {
		t.Fatalf("expected UpCount=2, got %d", snap.UpCount)
	}
	if snap.DownCount != 1 {
		t.Fatalf("expected DownCount=1, got %d", snap.DownCount)
	}
}

func TestSnapshot_IsIndependent(t *testing.T) {
	c := metrics.New()
	c.RecordScan()
	s1 := c.Snapshot()
	c.RecordScan()
	s2 := c.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Fatal("snapshots should be independent")
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := metrics.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.RecordScan()
			c.RecordAlert()
			c.RecordPortUp()
		}()
	}
	wg.Wait()
	snap := c.Snapshot()
	if snap.ScansTotal != 100 {
		t.Fatalf("expected 100 scans, got %d", snap.ScansTotal)
	}
}
