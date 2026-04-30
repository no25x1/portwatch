// Package audit maintains a bounded, thread-safe audit trail of port
// state-change events observed by portwatch.
//
// Each [Entry] records the host, port, previous state, current state,
// timestamp, and the source component that detected the change.
//
// Typical usage:
//
//	log := audit.New(500)
//	log.Record(audit.Entry{
//		Host:   "db.internal",
//		Port:   5432,
//		Prev:   "closed",
//		Curr:   "open",
//		Source: "poller",
//	})
//
//	// later — write NDJSON to stdout and reset
//	log.Flush(os.Stdout)
package audit
