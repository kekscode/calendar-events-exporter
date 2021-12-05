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

type IcalStore struct{}

type Event struct {
	Summary     string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	ID          string
}

// Returns a new event store of a given type
func NewEventStore(storeType string, targets []string) (EventStore, error) {

	switch storeType {
	case "ical":
		icsStore, err := newICalEventStore(targets)
		icsStore.Update()

		store := IcalStore{}
		return &store, err
	}
	return nil, errors.New("unknown store type")
}

func (st *IcalStore) Update() {}

func (st *IcalStore) GetEvents() []Event {
	return nil
}
