package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(host string, port int, open bool) scanner.PortEvent {
	return scanner.PortEvent{Host: host, Port: port, Open: open}
}

func TestOnlyOpen(t *testing.T) {
	events := []scanner.PortEvent{
		makeEvent("host1", 80, true),
		makeEvent("host1", 443, false),
		makeEvent("host2", 22, true),
	}
	result := filter.Apply(events, filter.OnlyOpen())
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestOnlyClosed(t *testing.T) {
	events := []scanner.PortEvent{
		makeEvent("host1", 80, true),
		makeEvent("host1", 443, false),
	}
	result := filter.Apply(events, filter.OnlyClosed())
	if len(result) != 1 || result[0].Port != 443 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestOnlyHosts(t *testing.T) {
	events := []scanner.PortEvent{
		makeEvent("alpha", 80, true),
		makeEvent("beta", 80, true),
		makeEvent("gamma", 80, true),
	}
	result := filter.Apply(events, filter.OnlyHosts("alpha", "gamma"))
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestOnlyPorts(t *testing.T) {
	events := []scanner.PortEvent{
		makeEvent("host1", 80, true),
		makeEvent("host1", 443, true),
		makeEvent("host1", 22, true),
	}
	result := filter.Apply(events, filter.OnlyPorts(80, 22))
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestChain(t *testing.T) {
	events := []scanner.PortEvent{
		makeEvent("host1", 80, true),
		makeEvent("host1", 443, false),
		makeEvent("host2", 80, true),
	}
	p := filter.Chain(filter.OnlyOpen(), filter.OnlyHosts("host1"))
	result := filter.Apply(events, p)
	if len(result) != 1 || result[0].Host != "host1" || result[0].Port != 80 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestApply_EmptyEvents(t *testing.T) {
	result := filter.Apply(nil, filter.OnlyOpen())
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %+v", result)
	}
}
