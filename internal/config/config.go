package config

import (
	"encoding/json"
	"os"
	"time"
)

// Host defines a single host and its ports to monitor.
type Host struct {
	Address string `json:"address"`
	Ports   []int  `json:"ports"`
}

// Config holds the full portwatch configuration.
type Config struct {
	Hosts    []Host        `json:"hosts"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	StateFile string       `json:"state_file"`
	LogLevel string        `json:"log_level"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  30 * time.Second,
		Timeout:   2 * time.Second,
		StateFile: "portwatch.state",
		LogLevel:  "info",
	}
}

// Load reads and parses a JSON config file from the given path.
// Missing optional fields fall back to defaults.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return cfg, err
	}

	if cfg.Interval <= 0 {
		cfg.Interval = DefaultConfig().Interval
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}
	if cfg.StateFile == "" {
		cfg.StateFile = DefaultConfig().StateFile
	}

	return cfg, nil
}
