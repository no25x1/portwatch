// Package filter provides composable predicates for filtering port scan events.
//
// Predicates can be combined using Chain to build AND-logic pipelines:
//
//	p := filter.Chain(
//		filter.OnlyOpen(),
//		filter.OnlyHosts("db1", "db2"),
//		filter.OnlyPorts(5432, 3306),
//	)
//	matched := filter.Apply(events, p)
package filter
