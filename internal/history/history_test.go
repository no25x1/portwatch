package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRecord_AppendsEntry(t *testing.T) {
	h := New(10)
	h.Record("localhost", 80, true)
	all := h.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if all[0].Host != "localhost" || all[0].Port != 80 || !all[0].Open {
		t.Errorf("unexpected entry: %+v", all[0])
	}
}

func TestRecord_Eviction(t *testing.T) {
	h := New(3)
	for i := 0; i < 5; i++ {
		h.Record("host", i, true)
	}
	all := h.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", len(all))
	}
	if all[0].Port != 2 {
		t.Errorf("expected oldest surviving port=2, got %d", all[0].Port)
	}
}

func TestSaveAndLoadJSON(t *testing.T) {
	h := New(10)
	h.Record("a.example.com", 443, true)
	h.Record("b.example.com", 22, false)

	tmp := filepath.Join(t.TempDir(), "history.json")
	if err := h.SaveJSON(tmp); err != nil {
		t.Fatalf("SaveJSON: %v", err)
	}

	h2 := New(10)
	if err := h2.LoadJSON(tmp); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	all := h2.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[1].Host != "b.example.com" || all[1].Port != 22 || all[1].Open {
		t.Errorf("unexpected second entry: %+v", all[1])
	}
}

func TestLoadJSON_InvalidFile(t *testing.T) {
	h := New(10)
	if err := h.LoadJSON("/nonexistent/path.json"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveJSON_ValidJSON(t *testing.T) {
	h := New(10)
	h.Record("localhost", 8080, true)
	tmp := filepath.Join(t.TempDir(), "out.json")
	_ = h.SaveJSON(tmp)
	data, _ := os.ReadFile(tmp)
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}
