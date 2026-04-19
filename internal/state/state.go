package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// PortKey uniquely identifies a host+port combination.
type PortKey struct {
	Host string
	Port int
}

// PortState records the last known state of a port.
type PortState struct {
	Open      bool
	LastSeen  time.Time
	LastCheck time.Time
}

// Store holds port states and can persist them to disk.
type Store struct {
	mu    sync.RWMutex
	states map[PortKey]PortState
	path  string
}

// New creates a new Store, loading persisted state from path if it exists.
func New(path string) (*Store, error) {
	s := &Store{
		states: make(map[PortKey]PortState),
		path:  path,
	}
	if path != "" {
		if err := s.load(); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return s, nil
}

// Update sets the current state for a port and returns whether it changed.
func (s *Store) Update(key PortKey, open bool, now time.Time) (changed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, exists := s.states[key]
	changed = !exists || prev.Open != open
	st := PortState{Open: open, LastCheck: now}
	if open {
		st.LastSeen = now
	} else if exists {
		st.LastSeen = prev.LastSeen
	}
	s.states[key] = st
	return changed
}

// Get returns the stored state for a key.
func (s *Store) Get(key PortKey) (PortState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.states[key]
	return st, ok
}

type persistEntry struct {
	Key   PortKey
	Value PortState
}

// Save persists state to disk.
func (s *Store) Save() error {
	if s.path == "" {
		return nil
	}
	s.mu.RLock()
	entries := make([]persistEntry, 0, len(s.states))
	for k, v := range s.states {
		entries = append(entries, persistEntry{Key: k, Value: v})
	}
	s.mu.RUnlock()
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entries)
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	var entries []persistEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return err
	}
	for _, e := range entries {
		s.states[e.Key] = e.Value
	}
	return nil
}
