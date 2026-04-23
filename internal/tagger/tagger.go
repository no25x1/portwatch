// Package tagger provides label-based tagging for scan targets,
// allowing events and results to carry arbitrary key-value metadata.
package tagger

import "sync"

// Tags is an immutable map of string key-value labels.
type Tags map[string]string

// Clone returns a shallow copy of the Tags map.
func (t Tags) Clone() Tags {
	out := make(Tags, len(t))
	for k, v := range t {
		out[k] = v
	}
	return out
}

// Has returns true if the tag key exists.
func (t Tags) Has(key string) bool {
	_, ok := t[key]
	return ok
}

// Get returns the value for a key and whether it was present.
func (t Tags) Get(key string) (string, bool) {
	v, ok := t[key]
	return v, ok
}

// Registry maps target identifiers (host:port) to their Tags.
type Registry struct {
	mu   sync.RWMutex
	tags map[string]Tags
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{tags: make(map[string]Tags)}
}

// Set associates tags with the given target key, replacing any existing entry.
func (r *Registry) Set(target string, t Tags) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[target] = t.Clone()
}

// Get returns a copy of the tags for target, and whether the target was found.
func (r *Registry) Get(target string) (Tags, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tags[target]
	if !ok {
		return nil, false
	}
	return t.Clone(), true
}

// Delete removes the tag entry for target.
func (r *Registry) Delete(target string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tags, target)
}

// All returns a snapshot of all registered targets and their tags.
func (r *Registry) All() map[string]Tags {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]Tags, len(r.tags))
	for k, v := range r.tags {
		out[k] = v.Clone()
	}
	return out
}
