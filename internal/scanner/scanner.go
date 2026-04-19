package scanner

import (
	"fmt"
	"net"
	"time"
)

// State represents the state of a port.
type State string

const (
	StateOpen   State = "open"
	StateClosed State = "closed"
)

// Result holds the result of a single port scan.
type Result struct {
	Host  string
	Port  int
	State State
	Error error
}

// Options configures scanner behaviour.
type Options struct {
	Timeout time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Timeout: 2 * time.Second,
	}
}

// CheckPort probes a single TCP port on host and returns a Result.
func CheckPort(host string, port int, opts Options) Result {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, opts.Timeout)
	if err != nil {
		return Result{Host: host, Port: port, State: StateClosed, Error: err}
	}
	conn.Close()
	return Result{Host: host, Port: port, State: StateOpen}
}

// CheckPorts probes multiple ports on a host concurrently and returns all results.
func CheckPorts(host string, ports []int, opts Options) []Result {
	results := make([]Result, len(ports))
	ch := make(chan Result, len(ports))

	for _, port := range ports {
		go func(p int) {
			ch <- CheckPort(host, p, opts)
		}(port)
	}

	for i := range ports {
		results[i] = <-ch
	}
	return results
}
