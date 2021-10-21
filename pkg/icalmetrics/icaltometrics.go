package icalmetrics 

import (
	"fmt"
	"net/http"

	ical "github.com/arran4/golang-ical"
)

type CalendarMetrics struct{}

func NewCalendarMetrics(calendarURL string) CalendarMetrics {
	c := CalendarMetrics{}
	resp, err := http.Get(calendarURL)
	if err != nil {
		fmt.Println("cannot fetch calendar: ", err)
		return c
	}
	defer resp.Body.Close()

	content, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		fmt.Println("cannot parse calendar data: ", err)
		return c
	}

	fmt.Printf("%v", content.Events())

	return c
}
