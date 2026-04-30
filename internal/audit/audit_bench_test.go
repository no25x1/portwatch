package audit

import (
	"io"
	"testing"
	"time"
)

func BenchmarkRecord(b *testing.B) {
	l := New(10000)
	e := Entry{
		Timestamp: time.Now().UTC(),
		Host:      "bench.host",
		Port:      8080,
		Prev:      "closed",
		Curr:      "open",
		Source:    "bench",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Record(e)
	}
}

func BenchmarkFlush(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := New(500)
		for j := 0; j < 500; j++ {
			l.Record(Entry{Host: "h", Port: j, Prev: "closed", Curr: "open", Source: "b"})
		}
		b.StartTimer()
		_ = l.Flush(io.Discard)
		b.StopTimer()
	}
}
