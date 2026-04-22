package watchdog_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/watchdog"
)

func BenchmarkWatchdog_Run(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := metrics.New()
		runner := &fakeRunner{}
		wd := watchdog.New(runner, m, watchdog.Options{
			CheckInterval: 1 * time.Millisecond,
			MaxErrors:     100,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		wd.Run(ctx)
		cancel()
	}
}
