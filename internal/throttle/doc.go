// Package throttle implements a fixed-window token-bucket rate limiter for
// controlling how many port scans portwatch dispatches within a given time
// window.
//
// Usage:
//
//	th := throttle.New(100, time.Second) // 100 scans per second
//	for _, target := range targets {
//		if err := th.Acquire(ctx); err != nil {
//			return err
//		}
//		go scan(target)
//	}
//
// A maxPerWindow of 0 (or negative) disables rate limiting entirely, allowing
// unbounded scan throughput.
package throttle
