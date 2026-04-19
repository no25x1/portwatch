package poller_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/poller"
)

func TestFromConfig_Valid(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.HostEntry{
			{Host: "example.com", Ports: []int{80, 443}},
			{Host: "db.internal", Ports: []int{5432}},
		},
	}
	targets, err := poller.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}
	if targets[0].Host != "example.com" || targets[0].Port != 80 {
		t.Errorf("unexpected first target: %+v", targets[0])
	}
}

func TestFromConfig_MissingHost(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.HostEntry{
			{Host: "", Ports: []int{80}},
		},
	}
	_, err := poller.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestFromConfig_InvalidPort(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.HostEntry{
			{Host: "example.com", Ports: []int{0}},
		},
	}
	_, err := poller.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestFromConfig_NoPorts(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.HostEntry{
			{Host: "example.com", Ports: []int{}},
		},
	}
	_, err := poller.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for no ports")
	}
}
