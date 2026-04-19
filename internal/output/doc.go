// Package output provides formatted event writing for portwatch.
//
// It supports two output formats:
//
//	- text: human-readable single-line entries suitable for terminal use
//	- json: newline-delimited JSON suitable for log aggregation pipelines
//
// Usage:
//
//	w := output.New(output.FormatJSON, os.Stdout)
//	w.Write(output.Event{
//		Host:      "10.0.0.1",
//		Port:      443,
//		State:     "open",
//		PrevState: "closed",
//		Timestamp: time.Now(),
//	})
package output
