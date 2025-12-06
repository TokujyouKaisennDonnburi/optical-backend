package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"

type CalendarQuery struct {
	calendarRepository repository.CalendarRepository
}

type EventQuery struct {
	eventRepository repository.EventRepository
}

func NewCalendarQuery(calendarRepository repository.CalendarRepository) *CalendarQuery {
	if calendarRepository == nil {
		panic("CalendarRepository is nil")
	}
	return &CalendarQuery{
		calendarRepository: calendarRepository,
	}
}

func NewEventQuery(eventRepository repository.EventRepository) *EventQuery {
	if eventRepository == nil {
		panic("eventRepository is nil")
	}
	return &EventQuery{
		eventRepository: eventRepository,
	}
}
