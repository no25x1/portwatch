// Package watchdog provides a self-healing supervisor for the portwatch
// polling runner.
//
// # Overview
//
// A Watchdog wraps any Runner (Start/Stop) and periodically inspects
// a metrics.Metrics snapshot.  When the number of consecutive scan
// errors reaches the configured MaxErrors threshold the watchdog calls
// Stop on the runner and launches a fresh goroutine that calls Start,
// allowing the poller to recover from transient network or DNS failures
// without operator intervention.
//
// # Usage
//
//	wd := watchdog.New(myPoller, m, watchdog.Options{
//	    CheckInterval: 30 * time.Second,
//	    MaxErrors:     10,
//	})
//	go wd.Run(ctx)
package watchdog
