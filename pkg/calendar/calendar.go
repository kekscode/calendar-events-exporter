package calendar

import (
	"net/http"

	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
)

type calendars struct {
	calendars []calendar
}

type calendar struct {
	URL      string
	calendar *ics.Calendar
	vevents  []*ics.VEvent
}

// NewCalendards returns a new Calendars struct
func newCalendars(targets []string) *calendars {
	c := calendars{}

	for _, target := range targets {
		c.calendars = append(c.calendars, calendar{
			URL:      target,
			calendar: ics.NewCalendarFor(target),
			vevents:  c.updateCalendars(),
		},
		)
	}
	return &c
}

// updateCalendars updates the Calendars struct with the latest content from the targets
func (c *calendars) updateCalendars() []*ics.VEvent {
	var vevents []*ics.VEvent

	for _, target := range c.calendars {
		resp, err := http.Get(target.URL)
		if err != nil {
			log.Println("cannot fetch calendar: ", err)
		}
		defer resp.Body.Close()

		target.calendar, err = ics.ParseCalendar(resp.Body)
		if err != nil {
			log.Println("cannot parse calendar data: ", err)
		}

		vevents = append(vevents, target.calendar.Events()...)
	}

	if len(vevents) > 0 {
		return vevents
	}

	return []*ics.VEvent{}
}
