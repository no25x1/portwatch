// Package backoff provides configurable exponential back-off helpers
// used when retrying failed port scans or alert deliveries.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Policy describes how delays are calculated between retry attempts.
type Policy struct {
	// Base is the initial delay before the first retry.
	Base time.Duration
	// Max caps the computed delay so it never grows unbounded.
	Max time.Duration
	// Factor is the multiplier applied after each attempt (e.g. 2.0).
	Factor float64
	// Jitter, when true, adds a random fraction of the computed delay to
	// spread load across concurrent workers.
	Jitter bool
}

// DefaultPolicy returns a Policy suitable for most portwatch retry sites.
func DefaultPolicy() Policy {
	return Policy{
		Base:   250 * time.Millisecond,
		Max:    30 * time.Second,
		Factor: 2.0,
		Jitter: true,
	}
}

// Delay returns the back-off duration for the given attempt number (0-based).
// attempt == 0 returns Base; each subsequent attempt multiplies by Factor up
// to Max. When Jitter is enabled, a random value in [0, delay) is added.
func (p Policy) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	delay := float64(p.Base) * math.Pow(p.Factor, float64(attempt))
	if max := float64(p.Max); delay > max {
		delay = max
	}

	if p.Jitter && delay > 0 {
		// Add up to 100 % of the computed delay as jitter.
		delay += rand.Float64() * delay //nolint:gosec // non-crypto use
		if max := float64(p.Max); delay > max {
			delay = max
		}
	}

	return time.Duration(delay)
}

// Sequence returns a channel that emits delays for attempts 0..n-1 and then
// closes. Callers can range over it to drive a retry loop with correct
// back-off timing without managing attempt counters themselves.
func (p Policy) Sequence(n int) <-chan time.Duration {
	ch := make(chan time.Duration, n)
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			ch <- p.Delay(i)
		}
	}()
	return ch
}
