package buffer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/buffer"
)

func TestConcurrentPush_NeverExceedsCap(t *testing.T) {
	const cap = 64
	const goroutines = 20
	const pushesEach = 50

	b := buffer.New(cap)
	var wg sync.WaitGroup

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < pushesEach; i++ {
				b.Push(buffer.Entry{
					Host:      "host",
					Port:      id*100 + i,
					Open:      i%2 == 0,
					Timestamp: time.Now(),
				})
			}
		}(g)
	}

	wg.Wait()

	if l := b.Len(); l > cap {
		t.Errorf("buffer length %d exceeds cap %d", l, cap)
	}
}

func TestDrainThenPush_BufferGrowsAgain(t *testing.T) {
	b := buffer.New(8)
	for i := 0; i < 8; i++ {
		b.Push(buffer.Entry{Host: "h", Port: i, Open: true, Timestamp: time.Now()})
	}

	drained := b.Drain()
	if len(drained) != 8 {
		t.Fatalf("expected 8 drained entries, got %d", len(drained))
	}
	if b.Len() != 0 {
		t.Fatalf("expected empty buffer after drain")
	}

	b.Push(buffer.Entry{Host: "new", Port: 9999, Open: false, Timestamp: time.Now()})
	if b.Len() != 1 {
		t.Errorf("expected len 1 after post-drain push, got %d", b.Len())
	}
	if all := b.All(); all[0].Host != "new" {
		t.Errorf("unexpected host after drain+push: %q", all[0].Host)
	}
}
