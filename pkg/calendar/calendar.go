package calendar

import (
	"fmt"
	"net/http"

	ics "github.com/arran4/golang-ical"
)

type Calendars struct {
	Calendars []Calendar
}

type Calendar struct {
	URL      string
	Calendar *ics.Calendar
}

// NewCalendards returns a new Calendars struct
func NewCalendars(targets []string) *[]Calendar {
	c := []Calendar{}

	for _, target := range targets {
		c = append(c, Calendar{
			URL:      target,
			Calendar: ics.NewCalendarFor(target),
		},
		)
	}
	return &c
}

// UpdateCalendars updates the Calendars struct with the latest content from the targets
// TODO: Use visitor pattern here like .Apply(func()...)
func (c *Calendars) UpdateCalendars() error {
	for _, target := range c.Calendars {
		resp, err := http.Get(target.URL)
		if err != nil {
			fmt.Println("cannot fetch calendar: ", err)
		}
		defer resp.Body.Close()

		target.Calendar, err = ics.ParseCalendar(resp.Body)
		if err != nil {
			fmt.Println("cannot parse calendar data: ", err)
		}
	}
	return nil
}

// Next step: Fetch and parse ical data
//
//resp, err := http.Get(targets[0])
//if err != nil {
//	fmt.Println("cannot fetch calendar: ", err)
//	return c
//}
//defer resp.Body.Close()

//content, err := ical.ParseCalendar(resp.Body)
//if err != nil {
//	fmt.Println("cannot parse calendar data: ", err)
//	return c
//}

//fmt.Printf("%v", content.Events())
