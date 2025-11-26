package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"

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
