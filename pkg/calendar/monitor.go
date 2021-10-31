package calendar

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
)

type Monitor struct {
	Events []ics.VEvent
}

func NewMonitor(targets []string) (*Monitor, error) {

	mon := Monitor{}
	cals := newCalendars(targets)
	for _, c := range *cals {
		fmt.Printf("Loading cal: %v", c)
	}

	return &mon, nil
}
