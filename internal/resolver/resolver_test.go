package resolver

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func mockLookup(addrs []string, err error, calls *atomic.Int32) func(context.Context, string) ([]string, error) {
	return func(_ context.Context, _ string) ([]string, error) {
		calls.Add(1)
		return addrs, err
	}
}

func TestResolve_ReturnsAddress(t *testing.T) {
	var calls atomic.Int32
	r := New(time.Minute)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	addr, err := r.Resolve(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "1.2.3.4" {
		t.Fatalf("want 1.2.3.4, got %s", addr)
	}
}

func TestResolve_CachesResult(t *testing.T) {
	var calls atomic.Int32
	r := New(time.Minute)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	for i := 0; i < 5; i++ {
		_, _ = r.Resolve(context.Background(), "example.com")
	}
	if calls.Load() != 1 {
		t.Fatalf("expected 1 lookup, got %d", calls.Load())
	}
}

func TestResolve_ExpiredEntryRefetched(t *testing.T) {
	var calls atomic.Int32
	r := New(10 * time.Millisecond)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	_, _ = r.Resolve(context.Background(), "host")
	time.Sleep(20 * time.Millisecond)
	_, _ = r.Resolve(context.Background(), "host")

	if calls.Load() != 2 {
		t.Fatalf("expected 2 lookups after expiry, got %d", calls.Load())
	}
}

func TestResolve_LookupError(t *testing.T) {
	r := New(time.Minute)
	r.lookup = mockLookup(nil, errors.New("dns failure"), nil)

	_, err := r.Resolve(context.Background(), "bad.host")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	var calls atomic.Int32
	r := New(time.Minute)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	_, _ = r.Resolve(context.Background(), "host")
	r.Invalidate("host")
	_, _ = r.Resolve(context.Background(), "host")

	if calls.Load() != 2 {
		t.Fatalf("expected 2 lookups after invalidation, got %d", calls.Load())
	}
}

func TestSize_ReflectsCache(t *testing.T) {
	var calls atomic.Int32
	r := New(time.Minute)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	_, _ = r.Resolve(context.Background(), "a")
	_, _ = r.Resolve(context.Background(), "b")

	if r.Size() != 2 {
		t.Fatalf("expected size 2, got %d", r.Size())
	}
	r.Invalidate("a")
	if r.Size() != 1 {
		t.Fatalf("expected size 1 after invalidation, got %d", r.Size())
	}
}
