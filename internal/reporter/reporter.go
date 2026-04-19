// Package reporter ties together summary, history, and output to produce
// periodic status reports for all monitored targets.
package reporter

import (
	"io"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/output"
	"github.com/yourorg/portwatch/internal/summary"
)

// Reporter renders a periodic digest of port-state activity.
type Reporter struct {
	builder *summary.Builder
	history *history.History
	writer  *output.Writer
}

// New creates a Reporter that writes to w using the given format ("text"/"json").
func New(w io.Writer, format string, h *history.History) *Reporter {
	return &Reporter{
		builder: summary.New(),
		history: h,
		writer:  output.New(w, format),
	}
}

// Record ingests a scanner event into the summary builder and history log.
func (r *Reporter) Record(host string, port int, open bool, ts time.Time) {
	r.builder.Upsert(host, port, open, ts)
	r.history.Record(host, port, open, ts)
}

// Flush builds the current summary report and writes it, then resets state.
func (r *Reporter) Flush() error {
	rpt := r.builder.Build()
	for _, entry := range rpt.Entries {
		if err := r.writer.Write(toEvent(entry)); err != nil {
			return err
		}
	}
	r.builder = summary.New()
	return nil
}

// toEvent converts a summary entry into the map shape output.Writer expects.
func toEvent(e summary.Entry) map[string]any {
	return map[string]any{
		"host":       e.Host,
		"port":       e.Port,
		"open":       e.Open,
		"last_seen":  e.LastSeen,
		"change_cnt": e.Changes,
	}
}
