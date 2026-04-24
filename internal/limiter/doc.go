// Package limiter provides a concurrency limiter for bounding the number of
// simultaneous operations within portwatch.
//
// # Overview
//
// Limiter wraps a buffered channel semaphore so callers can Acquire a slot
// before starting a scan and Release it when done. This prevents runaway
// goroutine creation when sweeping large host/port matrices.
//
// # Usage
//
//	l, err := limiter.New(16) // at most 16 concurrent scans
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := l.Acquire(ctx); err != nil {
//	    return err // context cancelled
//	}
//	defer l.Release()
//	// ... perform scan ...
package limiter
