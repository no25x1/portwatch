package poller

import (
	"fmt"

	"github.com/user/portwatch/internal/config"
)

// FromConfig converts a Config into a flat list of Targets.
func FromConfig(cfg *config.Config) ([]Target, error) {
	var targets []Target
	for _, h := range cfg.Hosts {
		if h.Host == "" {
			return nil, fmt.Errorf("host entry missing 'host' field")
		}
		if len(h.Ports) == 0 {
			return nil, fmt.Errorf("host %q has no ports defined", h.Host)
		}
		for _, p := range h.Ports {
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("invalid port %d for host %q", p, h.Host)
			}
			targets = append(targets, Target{Host: h.Host, Port: p})
		}
	}
	return targets, nil
}
