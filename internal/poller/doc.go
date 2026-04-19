// Package poller implements the core polling loop for portwatch.
//
// It periodically scans a list of Targets (host+port pairs) using the scanner
// package, records state changes via the state package, and triggers alerts
// through the alert package when a port transitions between open and closed.
//
// Typical usage:
//
//	cfg, _ := config.Load("portwatch.yaml")
//	targets, _ := poller.FromConfig(cfg)
//	st := state.New()
//	al := alert.New(nil)
//	p := poller.New(targets, cfg.Interval, st, al, scanner.DefaultOptions)
//	p.Run(ctx)
package poller
