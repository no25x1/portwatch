// Package filter provides predicates for selectively processing port events.
package filter

import "github.com/user/portwatch/internal/scanner"

// Predicate is a function that returns true if the event should be processed.
type Predicate func(e scanner.PortEvent) bool

// Chain combines multiple predicates with AND logic.
func Chain(predicates ...Predicate) Predicate {
	return func(e scanner.PortEvent) bool {
		for _, p := range predicates {
			if !p(e) {
				return false
			}
		}
		return true
	}
}

// OnlyOpen returns a predicate that passes only open port events.
func OnlyOpen() Predicate {
	return func(e scanner.PortEvent) bool {
		return e.Open
	}
}

// OnlyClosed returns a predicate that passes only closed port events.
func OnlyClosed() Predicate {
	return func(e scanner.PortEvent) bool {
		return !e.Open
	}
}

// OnlyHosts returns a predicate that passes events matching any of the given hosts.
func OnlyHosts(hosts ...string) Predicate {
	set := make(map[string]struct{}, len(hosts))
	for _, h := range hosts {
		set[h] = struct{}{}
	}
	return func(e scanner.PortEvent) bool {
		_, ok := set[e.Host]
		return ok
	}
}

// OnlyPorts returns a predicate that passes events matching any of the given ports.
func OnlyPorts(ports ...int) Predicate {
	set := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		set[p] = struct{}{}
	}
	return func(e scanner.PortEvent) bool {
		_, ok := set[e.Port]
		return ok
	}
}

// Apply returns only events that satisfy the predicate.
func Apply(events []scanner.PortEvent, p Predicate) []scanner.PortEvent {
	out := make([]scanner.PortEvent, 0, len(events))
	for _, e := range events {
		if p(e) {
			out = append(out, e)
		}
	}
	return out
}
