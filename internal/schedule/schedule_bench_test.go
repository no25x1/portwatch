package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

// BenchmarkScheduler_Throughput measures how quickly the scheduler can deliver
// ticks to a consumer under a very short interval.
func BenchmarkScheduler_Throughput(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := schedule.New(time.Millisecond)
	ch := s.Run(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-ch:
		case <-time.After(500 * time.Millisecond):
			b.Fatal("tick timeout")
		}
	}
}
