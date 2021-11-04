package calendar

import (
	ics "github.com/arran4/golang-ical"
)

// EventStore stores calendar events
type EventStore struct {
	// TODO: Secure this with a Write MUTEX lock
	// (defer unlock beachten)
	Events    []*ics.VEvent
	Calendars calendars
	targets   []string
}

// NewMonitor returns a new calendar monitor
func NewEventStore(targets []string) (*EventStore, error) {

	mon := EventStore{}
	mon.targets = targets

	return &mon, nil
}

func (m *EventStore) Update() {
	m.Calendars.updateCalendars()
	m.updateEvents()
}

func (m *EventStore) updateEvents() {
	cals := newCalendars(m.targets)
	cals.updateCalendars()
	m.Events = nil
	m.Events = append(m.Events, cals.vevents...)

	m.Calendars = *cals
}
