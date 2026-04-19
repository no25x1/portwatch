// Package reporter provides periodic digest reporting for portwatch.
//
// It combines the summary builder, history log, and output writer into a
// single surface so that callers only need to call Record on every scan
// event and Flush whenever they want a report emitted (e.g. on a ticker).
package reporter
