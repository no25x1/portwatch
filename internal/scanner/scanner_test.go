package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTCPServer starts a local TCP listener and returns its port and a stop func.
func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	// Accept connections in background to avoid connection resets.
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestCheckPort_Open(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	result := CheckPort("127.0.0.1", port, DefaultOptions())
	if result.State != StateOpen {
		t.Errorf("expected open, got %s (err: %v)", result.State, result.Error)
	}
}

func TestCheckPort_Closed(t *testing.T) {
	opts := Options{Timeout: 500 * time.Millisecond}
	// Port 1 is almost certainly closed in test environments.
	result := CheckPort("127.0.0.1", 1, opts)
	if result.State != StateClosed {
		t.Errorf("expected closed, got %s", result.State)
	}
}

func TestCheckPorts_Mixed(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	results := CheckPorts("127.0.0.1", []int{port, 1}, DefaultOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	openCount := 0
	for _, r := range results {
		if r.State == StateOpen {
			openCount++
		}
	}
	if openCount != 1 {
		t.Errorf("expected 1 open port, got %d", openCount)
	}
}
