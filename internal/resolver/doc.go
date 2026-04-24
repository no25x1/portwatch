// Package resolver provides a caching DNS resolver for portwatch.
//
// Hostnames supplied in configuration are resolved to IP addresses before
// each scan cycle. Results are cached with a configurable TTL so that
// repeated scans do not generate excessive DNS traffic.
//
// Usage:
//
//	r := resolver.New(30 * time.Second)
//	ip, err := r.Resolve(ctx, "example.com")
package resolver
