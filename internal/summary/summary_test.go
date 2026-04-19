package summary_test

import (
	"testing"
	"time"

	"portwatch/internal/summary"
)

func TestBuilder_UpsertAndBuild(t *testing.T) {
	b := summary.New()
	now := time.Now()

	b.Upsert("localhost", 80, true, now)
	b.Upsert("localhost", 443, false, now)
	b.Upsert("remotehost", 22, true, now)

	r := b.Build()

	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
	if r.TotalOpen() != 2 {
		t.Errorf("expected 2 open, got %d", r.TotalOpen())
	}
	if r.TotalClosed() != 1 {
		t.Errorf("expected 1 closed, got %d", r.TotalClosed())
	}
}

func TestBuilder_Upsert_Overwrites(t *testing.T) {
	b := summary.New()
	now := time.Now()

	b.Upsert("localhost", 80, true, now)
	b.Upsert("localhost", 80, false, now)

	r := b.Build()

	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", len(r.Entries))
	}
	if r.Entries[0].Open {
		t.Error("expected port to be closed after overwrite")
	}
}

func TestReport_EmptyTotals(t *testing.T) {
	b := summary.New()
	r := b.Build()

	if r.TotalOpen() != 0 || r.TotalClosed() != 0 {
		t.Error("expected zero totals for empty report")
	}
}

func TestKey(t *testing.T) {
	k := summary.Key("example.com", 8080)
	if k != "example.com:8080" {
		t.Errorf("unexpected key: %s", k)
	}
}
