// Package history provides a thread-safe, bounded event log for portwatch.
//
// Entries are recorded each time a port state change is detected. The log
// can be persisted to and restored from a JSON file, making it suitable for
// audit trails and post-incident review.
//
// When the log reaches its capacity, the oldest entries are evicted to make
// room for new ones, ensuring memory usage remains bounded over time.
//
// Usage:
//
//	h := history.New(500)
//	h.Record("web01", 443, true)
//	_ = h.SaveJSON("/var/lib/portwatch/history.json")
//
//	// Restore a previously saved log:
//	h2, err := history.LoadJSON("/var/lib/portwatch/history.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	entries := h2.Entries()
package history
