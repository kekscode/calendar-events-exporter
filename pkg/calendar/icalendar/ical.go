package icalendar

import (
	"net/http"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	ID          string
	Location    string
	Description string
	Summary     string
	StartTime   time.Time
	EndTime     time.Time
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
	VEvents   []*ics.VEvent
	Calendars calendars
	targets   []string
}

// newICSEventStore returns a new ical calendar event store
func NewEventStore(targets []string) (*ICSEventStore, error) {

	ical := ICSEventStore{}
	ical.targets = targets

	ical.Update()

	return &ical, nil
}

func (m *ICSEventStore) Events() []Event {
	iCalEvts := m.events()
	evts := []Event{}

	for _, iCalEvt := range iCalEvts {

		var startTime, endTime time.Time

		if iCalEvt.GetProperty(ics.ComponentPropertyDtStart) != nil {
			startTime, err := dateparse.ParseAny(iCalEvt.GetProperty(ics.ComponentPropertyDtStart).Value)
			if err != nil {
				log.Errorf("error with startTime: %v: %v", startTime, err)
			}
		}

		if iCalEvt.GetProperty(ics.ComponentPropertyDtEnd) != nil {
			endTime, err := dateparse.ParseAny(iCalEvt.GetProperty(ics.ComponentPropertyDtEnd).Value)
			if err != nil {
				log.Errorf("error with endTime: %v: %v", endTime, err)
			}
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

func (m *ICSEventStore) events() []*ics.VEvent {
	m.RLock()
	defer m.RUnlock()
	return m.VEvents
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
	m.VEvents = nil
	m.VEvents = append(m.VEvents, cals.vevents...)

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

		if resp == nil {
			log.Println("response is NIL!!!")
		}

		target.calendar, err = ics.ParseCalendar(resp.Body)
		if err != nil {
			// TODO: This panics if e.g. the calendar header is missing. Deal with this by skipping it.
			log.Println("cannot parse calendar data: ", err)
		}

		c.vevents = append(vevents, target.calendar.Events()...)
	}
}
