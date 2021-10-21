package calendar

import (
	"fmt"
	"net/http"

	ical "github.com/arran4/golang-ical"
)

type Calendar struct{}

func NewCalendar(calendarURL string) Calendar {
	c := Calendar{}
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

	fmt.Printf("%v", content)

	return c
}
