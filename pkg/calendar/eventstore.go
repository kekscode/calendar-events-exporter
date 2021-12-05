package calendar

import (
	"errors"
	"time"
)

// Generic store for calendar events
type EventStore interface {
	Update()
	GetEvents() []Event
}

type Event struct {
	Summary     string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	ID          string
}

// Returns a new event store of a given type
func NewEventStore(storeType string, targets []string) (*ICSEventStore, error) {

	switch storeType {
	case "ical":
		s, err := newICSEventStore(targets)
		return s, err
		//icsStore, err := newICSEventStore(targets)
		//icsStore.Update()

		//return &icsStore, err
	}
	return nil, errors.New("unknown store type")
}
