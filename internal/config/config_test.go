package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := map[string]any{
		"hosts": []map[string]any{
			{"address": "localhost", "ports": []int{80, 443}},
		},
		"interval": "10s",
		"timeout":  "1s",
		"state_file": "custom.state",
		"log_level": "debug",
	}
	path := writeTempConfig(t, raw)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 1 || cfg.Hosts[0].Address != "localhost" {
		t.Errorf("unexpected hosts: %+v", cfg.Hosts)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s interval, got %v", cfg.Interval)
	}
	if cfg.StateFile != "custom.state" {
		t.Errorf("expected custom.state, got %s", cfg.StateFile)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, map[string]any{})

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := DefaultConfig()
	if cfg.Interval != def.Interval {
		t.Errorf("expected default interval %v, got %v", def.Interval, cfg.Interval)
	}
	if cfg.StateFile != def.StateFile {
		t.Errorf("expected default state file %q, got %q", def.StateFile, cfg.StateFile)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
