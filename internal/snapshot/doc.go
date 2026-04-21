// Package snapshot provides point-in-time capture of port states and
// diff computation between consecutive snapshots.
//
// Typical usage:
//
//	b := snapshot.NewBuilder()
//	for _, r := range scanResults {
//		b.Record(r.Host, r.Port, r.Open)
//	}
//	next := b.Build()
//	changed := prev.Diff(next)
package snapshot
