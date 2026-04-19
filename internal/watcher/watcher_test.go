package watcher_test

import (
	"context"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/portwatch/internal/notify"
	"github.com/example/portwatch/internal/poller"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/state"
	"github.com/example/portwatch/internal/watcher"
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

func TestWatcher_DetectsOpenPort(t *testing.T) {
	port := freePort(t)
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	var called atomic.Int32
	ch := notify.ChannelFunc(func(_ context.Context, ev interface{}) error {
		called.Add(1)
		return nil
	})

	st := state.New()
	disp := notify.New(ch)
	targets := []poller.Target{{Host: "127.0.0.1", Ports: []int{port}}}

	w := watcher.New(watcher.Config{
		Targets:  targets,
		State:    st,
		Notifier: disp,
		Interval: 50 * time.Millisecond,
		Opts:     scanner.DefaultOptions(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if called.Load() == 0 {
		t.Error("expected at least one notification for open port")
	}
}

func TestWatcher_NoNotifyNoChange(t *testing.T) {
	port := freePort(t)

	var called atomic.Int32
	ch := notify.ChannelFunc(func(_ context.Context, ev interface{}) error {
		called.Add(1)
		return nil
	})

	st := state.New()
	disp := notify.New(ch)
	targets := []poller.Target{{Host: "127.0.0.1", Ports: []int{port}}}

	w := watcher.New(watcher.Config{
		Targets:  targets,
		State:    st,
		Notifier: disp,
		Interval: 40 * time.Millisecond,
		Opts:     scanner.DefaultOptions(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	// Port stays closed throughout; first scan fires a change, subsequent scans do not.
	if called.Load() > 1 {
		t.Errorf("expected at most 1 notification, got %d", called.Load())
	}
}
