// Package circuitbreaker provides a per-key circuit breaker for portwatch.
//
// A Breaker tracks consecutive scan failures for individual host:port targets.
// Once the configured failure threshold is reached the circuit opens and
// further Allow calls return ErrOpen, preventing wasted scan attempts against
// targets that are known to be unreachable.
//
// After the cooldown period the circuit moves to a half-open state: the next
// Allow call succeeds so the caller can attempt a probe. A successful probe
// (RecordSuccess) fully resets the breaker; another failure (RecordFailure)
// re-opens it for another cooldown cycle.
//
// Example:
//
//	br := circuitbreaker.New(3, 30*time.Second)
//	if err := br.Allow("192.168.1.1:22"); err != nil {
//		// skip scan — circuit is open
//		return
//	}
//	if err := scan(); err != nil {
//		br.RecordFailure("192.168.1.1:22")
//	} else {
//		br.RecordSuccess("192.168.1.1:22")
//	}
package circuitbreaker
