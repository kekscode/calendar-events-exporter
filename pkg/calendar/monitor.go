package calendar

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
)

// What's that about?
//kalender1,kalender2 <= prometheus_exporter:9310 => prometheus-metrics <= prometheus server (scrapen) <= Grafana WebUI (huebsche graphen)

// TODO: Add a "high level" event(s) store object and
// abstract away the ics.VEvents data structure

type Monitor struct {
	// TODO: Secure this with a Write MUTEX lock
	// (defer unlock beachten)
	Events    []*ics.VEvent
	Calendars calendars
	targets   []string
}

// NewMonitor returns a new calendar monitor
func NewMonitor(targets []string) (*Monitor, error) {

	mon := Monitor{}
	mon.targets = targets

	return &mon, nil
}

func (m *Monitor) Update() {
	m.Calendars.updateCalendars()
	m.updateEvents()
}

func (m *Monitor) updateEvents() {
	cals := newCalendars(m.targets)
	// FIXME: Not mockable
	// Better: Inject a monitor object to NewMonitor() to make it testable
	cals.updateCalendars()
	m.Events = nil

	for _, e := range cals.vevents {
		fmt.Printf("%v", e.GetProperty("SUMMARY"))
		m.Events = append(m.Events, e)
	}

	m.Calendars = *cals
}
