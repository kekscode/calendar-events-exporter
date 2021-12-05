package calendar

import (
	"net/http"
	"sync"

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

// ICalEventStore stores calendar events
type ICalEventStore struct {
	sync.RWMutex
	Events    []*ics.VEvent
	Calendars calendars
	targets   []string
}

// newICalEventStore returns a new ical calendar event store
func newICalEventStore(targets []string) (*ICalEventStore, error) {

	ical := ICalEventStore{}
	ical.targets = targets

	ical.Update()

	return &ical, nil
}

func (m *ICalEventStore) GetEvents() *[]Event {
	evts := []Event{}
	iCalEvts := m.getEvents()
	for _, iCalEvt := range iCalEvts {
		evts = append(evts, Event{
			Summary: iCalEvt.GetProperty(ics.ComponentPropertyUniqueId).Value,

			//TODO:
			//Location: iCalEvt.Location,
			//Start:    iCalEvt.Start.Time,
			//End:      iCalEvt.End.Time,
			//"uid":         store.Events[1].GetProperty(ics.ComponentPropertyUniqueId).Value,
			//"summary":     store.Events[1].GetProperty(ics.ComponentPropertySummary).Value,
			//"description": store.Events[1].GetProperty(ics.ComponentPropertyDescription).Value,
			//"location":    store.Events[1].GetProperty(ics.ComponentPropertyLocation).Value,
			//"dstart":      store.Events[1].GetProperty(ics.ComponentPropertyDtStart).Value,
			//"dend": store.Events[1].GetProperty(ics.ComponentPropertyDtEnd).Value,
		})
	}

	return &evts
}

func (m *ICalEventStore) getEvents() []*ics.VEvent {
	m.RLock()
	defer m.RUnlock()
	return m.Events
}

// Updates calendar events in the store
func (m *ICalEventStore) Update() {
	m.Lock()
	defer m.Unlock()

	m.Calendars.updateCalendars()
	m.updateEvents()
}

func (m *ICalEventStore) updateEvents() {
	cals := newCalendars(m.targets)
	cals.updateCalendars()
	m.Events = nil
	m.Events = append(m.Events, cals.vevents...)

	m.Calendars = *cals
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
			// TODO: This panics if e.g. the calendar header is missing. Deal with this by skipping it.
			log.Println("cannot parse calendar data: ", err)
		}

		c.vevents = append(vevents, target.calendar.Events()...)
	}
}
