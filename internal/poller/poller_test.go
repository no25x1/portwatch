package poller_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/poller"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestPoller_DetectsOpen(t *testing.T) {
	port := freePort(t)
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	st := state.New()
	al := alert.New(nil)
	targets := []poller.Target{{Host: "127.0.0.1", Port: port}}
	opts := scanner.Options{Timeout: 500 * time.Millisecond}

	p := poller.New(targets, 50*time.Millisecond, st, al, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	p.Run(ctx) //nolint:errcheck

	result, ok := st.Get("127.0.0.1", port)
	if !ok {
		t.Fatal("expected state entry")
	}
	if !result.Open {
		t.Error("expected port to be open")
	}
}

func TestPoller_DetectsClosed(t *testing.T) {
	port := freePort(t)

	st := state.New()
	al := alert.New(nil)
	targets := []poller.Target{{Host: "127.0.0.1", Port: port}}
	opts := scanner.Options{Timeout: 200 * time.Millisecond}

	p := poller.New(targets, 50*time.Millisecond, st, al, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	p.Run(ctx) //nolint:errcheck

	result, ok := st.Get("127.0.0.1", port)
	if !ok {
		t.Fatal("expected state entry")
	}
	if result.Open {
		t.Error("expected port to be closed")
	}
}
