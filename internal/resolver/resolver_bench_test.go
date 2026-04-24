package resolver

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkResolve_CacheHit(b *testing.B) {
	var calls atomic.Int32
	r := New(time.Hour)
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	// Prime the cache.
	_, _ = r.Resolve(context.Background(), "bench.host")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = r.Resolve(context.Background(), "bench.host")
		}
	})
}

func BenchmarkResolve_CacheMiss(b *testing.B) {
	var calls atomic.Int32
	r := New(time.Nanosecond) // expire immediately to force miss every time
	r.lookup = mockLookup([]string{"1.2.3.4"}, nil, &calls)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Resolve(context.Background(), "bench.host")
	}
}
