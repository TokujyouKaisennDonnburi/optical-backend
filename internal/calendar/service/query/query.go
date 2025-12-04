package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
)

// Event
type EventQuery struct {
	eventRepository repository.EventRepository
}

func NewEventQuery(eventRepo repository.EventRepository) *EventQuery {
	return &EventQuery{
		eventRepository: eventRepo,
	}
}

type CalendarQuery struct {
	calendarRepository repository.CalendarRepository
}

func NewCalendarQuery(calendarRepository repository.CalendarRepository) *CalendarQuery {
	if calendarRepository == nil {
		panic("CalendarRepository is nil")
	}
	return &CalendarQuery{
		calendarRepository: calendarRepository,
	}
}
