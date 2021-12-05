package calendar

import (
	"net/http"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
)

type IcalStore struct{}

func (st *IcalStore) Update() {}

func (st *IcalStore) GetEvents() []Event {
	return nil
}

type calendars struct {
	calendars []calendar
	vevents   []*ics.VEvent
}

type calendar struct {
	URL      string
	calendar *ics.Calendar
}

// ICSEventStore stores calendar events
type ICSEventStore struct {
	sync.RWMutex
	Events    []*ics.VEvent
	Calendars calendars
	targets   []string
}

// newICSEventStore returns a new ical calendar event store
func newICSEventStore(targets []string) (*ICSEventStore, error) {

	ical := ICSEventStore{}
	ical.targets = targets

	ical.Update()

	return &ical, nil
}

func (m *ICSEventStore) GetEvents() []Event {
	evts := []Event{}
	iCalEvts := m.getEvents()
	for _, iCalEvt := range iCalEvts {

		startTime, err := time.Parse("20210021T175157Z", iCalEvt.GetProperty(ics.ComponentPropertyDtStart).Value)
		if err != nil {
			log.Errorf("error: %v", err)
		}

		endTime, err := time.Parse("20210021T175157Z", iCalEvt.GetProperty(ics.ComponentPropertyDtEnd).Value)
		if err != nil {
			log.Errorf("error: %v", err)
		}

		evts = append(evts, Event{
			ID:          iCalEvt.GetProperty(ics.ComponentPropertyUniqueId).Value,
			Location:    iCalEvt.GetProperty(ics.ComponentPropertyLocation).Value,
			Summary:     iCalEvt.GetProperty(ics.ComponentPropertySummary).Value,
			Description: iCalEvt.GetProperty(ics.ComponentPropertyDescription).Value,
			StartTime:   startTime,
			EndTime:     endTime,
		})
	}

	return evts
}

func (m *ICSEventStore) getEvents() []*ics.VEvent {
	m.RLock()
	defer m.RUnlock()
	return m.Events
}

// Updates calendar events in the store
func (m *ICSEventStore) Update() {
	m.Lock()
	defer m.Unlock()

	m.Calendars.updateCalendars()
	m.updateEvents()
}

func (m *ICSEventStore) updateEvents() {
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
