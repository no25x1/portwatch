// Package semaphore implements a simple counting semaphore used by portwatch
// to cap the number of concurrent outbound TCP probes.
//
// Usage:
//
//	sem, err := semaphore.New(32)
//	if err != nil { ... }
//
//	if err := sem.Acquire(ctx); err != nil {
//		// context cancelled
//	}
//	defer sem.Release()
//	// ... perform scan ...
package semaphore
