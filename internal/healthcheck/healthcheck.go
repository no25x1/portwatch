// Package healthcheck provides an HTTP endpoint exposing runtime health
// and metrics for portwatch.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the current health of the portwatch process.
type Status struct {
	OK        bool              `json:"ok"`
	Uptime    string            `json:"uptime"`
	CheckedAt time.Time         `json:"checked_at"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Server exposes a /healthz HTTP endpoint.
type Server struct {
	mu      sync.RWMutex
	started time.Time
	meta    map[string]string
	mux     *http.ServeMux
}

// New creates a new Server ready to serve health status.
func New() *Server {
	s := &Server{
		started: time.Now(),
		meta:    make(map[string]string),
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("/healthz", s.handleHealth)
	return s
}

// SetMeta attaches an arbitrary key/value pair to the health response.
func (s *Server) SetMeta(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.meta[key] = value
}

// Handler returns the underlying http.Handler so it can be mounted externally.
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	metaCopy := make(map[string]string, len(s.meta))
	for k, v := range s.meta {
		metaCopy[k] = v
	}
	s.mu.RUnlock()

	status := Status{
		OK:        true,
		Uptime:    time.Since(s.started).Round(time.Second).String(),
		CheckedAt: time.Now().UTC(),
		Meta:      metaCopy,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}
