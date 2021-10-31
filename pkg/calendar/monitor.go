package calendar

import (
	ics "github.com/arran4/golang-ical"
)

type Monitor struct {
	Events []*ics.VEvent
}

// NewMonitor returns a new calendar monitor
func NewMonitor(targets []string) (*Monitor, error) {

	mon := Monitor{}
	cals := newCalendars(targets)

	cals.updateCalendars()

	for _, e := range cals.calendars {
		mon.Events = e.calendar.Events()
	}

	return &mon, nil
}
