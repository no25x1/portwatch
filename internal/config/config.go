// Package config loads portwatch configuration from a TOML file.
package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Target describes a host and the ports to monitor on it.
type Target struct {
	Host  string `toml:"host"`
	Ports []int  `toml:"ports"`
}

// Config holds the full portwatch configuration.
type Config struct {
	Interval  int      `toml:"interval_seconds"`
	StateFile string   `toml:"state_file"`
	LogFormat string   `toml:"log_format"`
	Targets   []Target `toml:"targets"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  30,
		StateFile: "portwatch.state",
		LogFormat: "text",
	}
}

// Load reads a TOML config file from path, applying defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, err
	}
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultConfig().Interval
	}
	if cfg.StateFile == "" {
		cfg.StateFile = DefaultConfig().StateFile
	}
	if cfg.LogFormat == "" {
		cfg.LogFormat = DefaultConfig().LogFormat
	}
	return cfg, nil
}
