package calendar

import (
	"fmt"
	"net/http"

	ics "github.com/arran4/golang-ical"
)

type calendars struct {
	calendars []calendar
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
func (c *calendars) updateCalendars() error {
	for _, target := range c.calendars {
		resp, err := http.Get(target.URL)
		if err != nil {
			fmt.Println("cannot fetch calendar: ", err)
		}
		defer resp.Body.Close()

		target.calendar, err = ics.ParseCalendar(resp.Body)
		if err != nil {
			fmt.Println("cannot parse calendar data: ", err)
		}
	}
	return nil
}
