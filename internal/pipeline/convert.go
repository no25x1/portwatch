package pipeline

import (
	"github.com/yourorg/portwatch/internal/filter"
	"github.com/yourorg/portwatch/internal/notify"
)

// toFilterEvent converts a pipeline Event to the type expected by filter.Func.
func toFilterEvent(e Event) filter.Event {
	return filter.Event{
		Host: e.Host,
		Port: e.Port,
		Open: e.Open,
	}
}

// toNotifyEvent converts a pipeline Event to the type expected by
// notify.Dispatcher.Dispatch.
func toNotifyEvent(e Event) notify.Event {
	return notify.Event{
		Host: e.Host,
		Port: e.Port,
		Open: e.Open,
	}
}
