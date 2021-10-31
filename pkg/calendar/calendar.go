package calendar

import (
	"net/http"

	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
)

type calendars struct {
	calendars []calendar
	vevents   []*ics.VEvent
}

type calendar struct {
	URL      string
	calendar *ics.Calendar
}

// NewCalendards returns a new Calendars struct
func newCalendars(targets []string) *calendars {
	c := calendars{}

	for _, target := range targets {
		c.calendars = append(c.calendars, calendar{
			URL:      target,
			calendar: ics.NewCalendarFor(target),
		},
		)
	}
	return &c
}

// updateCalendars updates the Calendars struct with the latest content from the targets
func (c *calendars) updateCalendars() {
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

		c.vevents = append(vevents, target.calendar.Events()...)
	}
}
