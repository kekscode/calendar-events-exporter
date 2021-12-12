package calendar

import (
	"fmt"
	"time"

	"github.com/kekscode/calendar-events-exporter/pkg/calendar/icalendar"
)

type Event struct {
	ID          string
	Location    string
	Description string
	Summary     string
	StartTime   time.Time
	EndTime     time.Time
}

// Generic store for calendar events
type EventStore interface {
	Update()
	Events() []Event
}

// Returns a new event store of a given type
func NewEventStore(storeType string, targets []string) (EventStore, error) {

	switch storeType {
	case "ical":
		s, err := icalendar.NewEventStore(targets)
		return s, err
	}
	return nil, fmt.Errorf("unknown store type \"%v\"", storeType)
}
