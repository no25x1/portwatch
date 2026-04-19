package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
)

func makeEvent() notify.Event {
	return notify.Event{
		Host:       "localhost",
		Port:       8080,
		PrevState:  "closed",
		CurrState:  "open",
		OccurredAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestDispatch_StdoutChannel(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdoutChannel(&buf)
	d := notify.New(ch)

	errs := d.Dispatch(makeEvent())
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port in output, got: %s", out)
	}
	if !strings.Contains(out, "closed -> open") {
		t.Errorf("expected state transition in output, got: %s", out)
	}
}

type failChannel struct{}

func (f *failChannel) Send(_ notify.Event) error {
	return errors.New("send failed")
}

func TestDispatch_CollectsErrors(t *testing.T) {
	d := notify.New(&failChannel{}, &failChannel{})
	errs := d.Dispatch(makeEvent())
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestDispatch_PartialFailure(t *testing.T) {
	var buf bytes.Buffer
	good := notify.NewStdoutChannel(&buf)
	d := notify.New(good, &failChannel{})

	errs := d.Dispatch(makeEvent())
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if buf.Len() == 0 {
		t.Error("expected good channel to have written output")
	}
}

func TestNewStdoutChannel_DefaultsToStdout(t *testing.T) {
	ch := notify.NewStdoutChannel(nil)
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}
