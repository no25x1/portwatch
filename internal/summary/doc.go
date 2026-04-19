// Package summary aggregates port scan results from the poller
// into structured reports. A Builder accumulates PortStatus entries
// via Upsert and produces a Report snapshot via Build.
//
// Reports expose convenience methods such as TotalOpen and TotalClosed
// and are suitable for rendering by the output package or exporting
// as JSON for downstream tooling.
package summary
