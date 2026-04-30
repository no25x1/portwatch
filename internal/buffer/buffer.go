// Package buffer provides a fixed-capacity ring buffer for port scan events.
// Older entries are evicted automatically once the buffer reaches capacity.
package buffer

import (
	"sync"
	"time"
)

// Entry holds a single buffered scan event.
type Entry struct {
	Host      string
	Port      int
	Open      bool
	Timestamp time.Time
}

// Buffer is a thread-safe, fixed-capacity ring buffer.
type Buffer struct {
	mu       sync.Mutex
	entries  []Entry
	cap      int
	head     int
	size     int
}

// New creates a Buffer with the given capacity. Panics if cap < 1.
func New(cap int) *Buffer {
	if cap < 1 {
		panic("buffer: capacity must be at least 1")
	}
	return &Buffer{
		entries: make([]Entry, cap),
		cap:     cap,
	}
}

// Push appends an entry, evicting the oldest entry when the buffer is full.
func (b *Buffer) Push(e Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := (b.head + b.size) % b.cap
	b.entries[idx] = e
	if b.size < b.cap {
		b.size++
	} else {
		// Overwrite oldest — advance head.
		b.head = (b.head + 1) % b.cap
	}
}

// All returns a snapshot of all buffered entries in insertion order.
func (b *Buffer) All() []Entry {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]Entry, b.size)
	for i := 0; i < b.size; i++ {
		out[i] = b.entries[(b.head+i)%b.cap]
	}
	return out
}

// Len returns the current number of entries in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.size
}

// Drain returns all entries and resets the buffer to empty.
func (b *Buffer) Drain() []Entry {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]Entry, b.size)
	for i := 0; i < b.size; i++ {
		out[i] = b.entries[(b.head+i)%b.cap]
	}
	b.head = 0
	b.size = 0
	return out
}
