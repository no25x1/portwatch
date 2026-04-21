// Package reporter provides periodic digest reporting for portwatch.
//
// It combines the summary builder, history log, and output writer into a
// single surface so that callers only need to call Record on every scan
// event and Flush whenever they want a report emitted (e.g. on a ticker).
//
// Typical usage:
//
//	r, err := reporter.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer r.Close()
//
//	// In the scan loop:
//	r.Record(event)
//
//	// On a time.Ticker or similar:
//	if err := r.Flush(); err != nil {
//		log.Printf("flush: %v", err)
//	}
package reporter
