// Package eventlog implements a bounded, thread-safe append-only log for
// port-state-change events produced by the portwatch scanning pipeline.
//
// Typical usage:
//
//	log := eventlog.New(500)
//	log.Append(eventlog.Entry{
//		Timestamp: time.Now(),
//		Host:      "10.0.0.1",
//		Port:      443,
//		State:     "open",
//	})
//	entries := log.Drain() // consume and clear
//
// The log evicts the oldest entry once the configured capacity is reached,
// ensuring memory usage stays bounded during long-running scans.
package eventlog
