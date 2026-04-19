package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/output"
)

func makeEvent() output.Event {
	return output.Event{
		Host:      "localhost",
		Port:      8080,
		State:     "open",
		PrevState: "closed",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatText, &buf)
	e := makeEvent()
	if err := w.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "localhost:8080") {
		t.Errorf("expected host:port in output, got: %s", got)
	}
	if !strings.Contains(got, "closed -> open") {
		t.Errorf("expected state transition in output, got: %s", got)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatJSON, &buf)
	e := makeEvent()
	if err := w.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded output.Event
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Host != e.Host || decoded.Port != e.Port {
		t.Errorf("decoded event mismatch: %+v", decoded)
	}
}

func TestNew_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	w := output.New("", &buf)
	e := makeEvent()
	if err := w.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.HasPrefix(buf.String(), "{") {
		t.Error("expected text format, got JSON")
	}
}
