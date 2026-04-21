package snapshot_test

import (
	"testing"

	"github.com/your-org/portwatch/internal/snapshot"
)

func buildSnapshot(records [][3]interface{}) *snapshot.Snapshot {
	b := snapshot.NewBuilder()
	for _, r := range records {
		b.Record(r[0].(string), r[1].(int), r[2].(bool))
	}
	return b.Build()
}

func TestKey(t *testing.T) {
	if got := snapshot.Key("localhost", 8080); got != "localhost:8080" {
		t.Fatalf("expected localhost:8080, got %s", got)
	}
}

func TestBuilder_RecordAndGet(t *testing.T) {
	s := buildSnapshot([][3]interface{}{
		{"host1", 80, true},
		{"host1", 443, false},
	})

	e, ok := s.Get("host1", 80)
	if !ok || !e.Open {
		t.Fatal("expected host1:80 to be open")
	}
	e, ok = s.Get("host1", 443)
	if !ok || e.Open {
		t.Fatal("expected host1:443 to be closed")
	}
	_, ok = s.Get("host2", 80)
	if ok {
		t.Fatal("expected missing entry to return false")
	}
}

func TestBuilder_All_Sorted(t *testing.T) {
	s := buildSnapshot([][3]interface{}{
		{"zhost", 80, true},
		{"ahost", 80, true},
		{"mhost", 80, false},
	})
	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	if all[0].Host != "ahost" || all[1].Host != "mhost" || all[2].Host != "zhost" {
		t.Fatalf("entries not sorted: %v", all)
	}
}

func TestDiff_DetectsStateChange(t *testing.T) {
	prev := buildSnapshot([][3]interface{}{
		{"host1", 80, true},
		{"host1", 443, false},
	})
	next := buildSnapshot([][3]interface{}{
		{"host1", 80, false}, // changed
		{"host1", 443, false}, // unchanged
	})
	changed := prev.Diff(next)
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].Port != 80 || changed[0].Open {
		t.Fatalf("unexpected changed entry: %+v", changed[0])
	}
}

func TestDiff_NewEntryIsChange(t *testing.T) {
	prev := buildSnapshot([][3]interface{}{
		{"host1", 80, true},
	})
	next := buildSnapshot([][3]interface{}{
		{"host1", 80, true},
		{"host1", 8080, true}, // newly appeared
	})
	changed := prev.Diff(next)
	if len(changed) != 1 || changed[0].Port != 8080 {
		t.Fatalf("expected new port 8080 to appear in diff, got %v", changed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	prev := buildSnapshot([][3]interface{}{{"host1", 80, true}})
	next := buildSnapshot([][3]interface{}{{"host1", 80, true}})
	if diff := prev.Diff(next); len(diff) != 0 {
		t.Fatalf("expected empty diff, got %v", diff)
	}
}

func TestSnapshot_CapturedAt_NotZero(t *testing.T) {
	s := snapshot.NewBuilder().Build()
	if s.CapturedAt().IsZero() {
		t.Fatal("CapturedAt should not be zero")
	}
}
