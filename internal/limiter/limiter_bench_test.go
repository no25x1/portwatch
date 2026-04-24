package limiter

import (
	"context"
	"testing"
)

// BenchmarkAcquireRelease measures the overhead of a single Acquire/Release
// cycle on an uncontended limiter.
func BenchmarkAcquireRelease(b *testing.B) {
	l, _ := New(b.N + 1)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Acquire(ctx)
		l.Release()
	}
}

// BenchmarkAcquireRelease_Contended measures throughput when goroutines compete
// for a small pool of slots.
func BenchmarkAcquireRelease_Contended(b *testing.B) {
	l, _ := New(4)
	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = l.Acquire(ctx)
			l.Release()
		}
	})
}
