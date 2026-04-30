// Package buffer implements a thread-safe fixed-capacity ring buffer for
// portwatch scan events. When the buffer is full, the oldest entry is
// silently evicted to make room for the newest one.
//
// Typical usage:
//
//	buf := buffer.New(512)
//	buf.Push(buffer.Entry{Host: "10.0.0.1", Port: 443, Open: true, Timestamp: time.Now()})
//	entries := buf.All()
package buffer
