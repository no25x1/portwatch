package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeEvent(prev, curr bool) Event {
	return Event{
		Host:      "localhost",
		Port:      8080,
		PrevOpen:  prev,
		CurrOpen:  curr,
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestNotify_PortCameUp(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Notify(makeEvent(false, true))
	out := buf.String()
	if !strings.Contains(out, string(LevelInfo)) {
		t.Errorf("expected INFO level, got: %s", out)
	}
	if !strings.Contains(out, "localhost:8080") {
		t.Errorf("expected host:port in output, got: %s", out)
	}
}

func TestNotify_PortWentDown(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Notify(makeEvent(true, false))
	out := buf.String()
	if !strings.Contains(out, string(LevelError)) {
		t.Errorf("expected ERROR level, got: %s", out)
	}
}

func TestNotify_NoChange(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Notify(makeEvent(true, true))
	out := buf.String()
	if !strings.Contains(out, string(LevelWarn)) {
		t.Errorf("expected WARN level, got: %s", out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}
