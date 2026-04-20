package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleHealth_ReturnsOK(t *testing.T) {
	s := New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var st Status
	if err := json.NewDecoder(rec.Body).Decode(&st); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !st.OK {
		t.Error("expected ok=true")
	}
	if st.CheckedAt.IsZero() {
		t.Error("expected non-zero checked_at")
	}
}

func TestHandleHealth_UptimePresent(t *testing.T) {
	s := New()
	time.Sleep(10 * time.Millisecond)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Handler().ServeHTTP(rec, req)

	var st Status
	_ = json.NewDecoder(rec.Body).Decode(&st)
	if st.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestSetMeta_AppearsInResponse(t *testing.T) {
	s := New()
	s.SetMeta("version", "1.2.3")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Handler().ServeHTTP(rec, req)

	var st Status
	_ = json.NewDecoder(rec.Body).Decode(&st)
	if st.Meta["version"] != "1.2.3" {
		t.Errorf("expected version=1.2.3, got %q", st.Meta["version"])
	}
}

func TestNew_DefaultMeta_Empty(t *testing.T) {
	s := New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Handler().ServeHTTP(rec, req)

	var st Status
	_ = json.NewDecoder(rec.Body).Decode(&st)
	if len(st.Meta) != 0 {
		t.Errorf("expected empty meta, got %v", st.Meta)
	}
}
