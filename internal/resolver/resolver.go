// Package resolver translates hostnames to IP addresses and caches
// the results for a configurable TTL to reduce DNS lookup overhead.
package resolver

import (
	"context"
	"net"
	"sync"
	"time"
)

// Entry holds a resolved address and its expiry time.
type Entry struct {
	Addr    string
	Expires time.Time
}

// Resolver caches DNS lookups with a TTL.
type Resolver struct {
	mu    sync.RWMutex
	cache map[string]Entry
	ttl   time.Duration
	lookup func(ctx context.Context, host string) ([]string, error)
}

// New returns a Resolver with the given TTL.
func New(ttl time.Duration) *Resolver {
	return &Resolver{
		cache:  make(map[string]Entry),
		ttl:    ttl,
		lookup: defaultLookup,
	}
}

func defaultLookup(ctx context.Context, host string) ([]string, error) {
	return net.DefaultResolver.LookupHost(ctx, host)
}

// Resolve returns the first IP address for host, using a cached result
// if one exists and has not expired.
func (r *Resolver) Resolve(ctx context.Context, host string) (string, error) {
	r.mu.RLock()
	entry, ok := r.cache[host]
	r.mu.RUnlock()

	if ok && time.Now().Before(entry.Expires) {
		return entry.Addr, nil
	}

	addrs, err := r.lookup(ctx, host)
	if err != nil {
		return "", err
	}

	addr := addrs[0]
	r.mu.Lock()
	r.cache[host] = Entry{Addr: addr, Expires: time.Now().Add(r.ttl)}
	r.mu.Unlock()

	return addr, nil
}

// Invalidate removes a cached entry for host.
func (r *Resolver) Invalidate(host string) {
	r.mu.Lock()
	delete(r.cache, host)
	r.mu.Unlock()
}

// Size returns the number of entries currently in the cache.
func (r *Resolver) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.cache)
}
