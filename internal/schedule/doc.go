// Package schedule provides a Scheduler that emits ticks at a configurable
// interval. The first tick fires immediately so that callers perform work
// on startup without waiting for the first interval to elapse.
//
// Usage:
//
//	s := schedule.New(30 * time.Second)
//	for tick := range s.Run(ctx) {
//		_ = tick // perform periodic work
//	}
package schedule
