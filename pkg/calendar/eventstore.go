package calendar

import (
	"fmt"
	"time"
)

// EventStore stores and updates calendar events
type EventStore interface {
	Update()
	Events() []Event
}

// Generic calendar event datatype
type Event struct {
	ID          string
	Location    string
	Description string
	Summary     string
	StartTime   time.Time
	EndTime     time.Time
}

// NewEventStore creates a new event store for a given type and a list of data source targets
func NewEventStore(storeType string, targets []string) (EventStore, error) {
	switch storeType {
	case "ical":
		s, err := NewICSEventStore(targets)
		s.Update()
		return s, err
	}
	return nil, fmt.Errorf("unknown store type \"%v\"", storeType)
}
