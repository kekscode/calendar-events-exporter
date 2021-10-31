package calendar

import (
	"fmt"

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
		fmt.Printf("%v", e.calendar.Events())
		mon.Events = append(mon.Events, e.calendar.Events()...)
	}

	return &mon, nil
}
