package resolver_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/resolver"
)

// TestConcurrentResolve verifies the cache is safe under concurrent access.
func TestConcurrentResolve(t *testing.T) {
	r := resolver.New(time.Minute)

	// Override the internal lookup via a known-good loopback address so the
	// test does not require external DNS.
	// We exercise the exported API only; the race detector validates safety.
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// localhost is always resolvable; ignore errors on restricted envs.
			_, _ = r.Resolve(context.Background(), "localhost")
		}()
	}
	wg.Wait()

	if r.Size() > 1 {
		t.Fatalf("expected at most 1 cache entry for localhost, got %d", r.Size())
	}
}

// TestInvalidate_ConcurrentSafe invalidates while resolvers are running.
func TestInvalidate_ConcurrentSafe(t *testing.T) {
	r := resolver.New(time.Minute)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = r.Resolve(context.Background(), "localhost")
			r.Invalidate("localhost")
		}()
	}
	wg.Wait()
	// No panic or race = pass.
}
