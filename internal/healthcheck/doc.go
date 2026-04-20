// Package healthcheck exposes a lightweight HTTP /healthz endpoint for the
// portwatch process. It reports liveness, uptime, and arbitrary metadata
// (such as version or active target count) as a JSON payload.
//
// Usage:
//
//	s := healthcheck.New()
//	s.SetMeta("version", buildVersion)
//	go s.ListenAndServe(ctx, ":9090")
package healthcheck
