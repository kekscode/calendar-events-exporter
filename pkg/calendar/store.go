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
func NewEventStore(storeType string, targets []string) (*EventStore, error) {
	switch storeType {
	case "ical":
		store, err := newICalEventStore(targets)
		return store, err
	}
	return nil, errors.New("unknown store type")
}

func Update(st *EventStore) {
}
func GetEvents(st *EventStore) []Event {
	return nil
}
