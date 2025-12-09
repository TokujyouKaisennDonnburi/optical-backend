package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
)

type CalendarHttpHandler struct {
	eventCommand    *command.EventCommand
	eventQuery      *query.EventQuery
	calendarCommand *command.CalendarCommand
	calendarQuery   *query.CalendarQuery
}

func NewCalendarHttpHandler(
	eventCommand *command.EventCommand,
	eventQuery *query.EventQuery,
	calendarCommand *command.CalendarCommand,
	calendarQuery *query.CalendarQuery,
) *CalendarHttpHandler {
	if eventCommand == nil {
		panic("EventCommand is nil")
	}
	if eventQuery == nil {
		panic("EventQuery is nil")
	}
	if calendarCommand == nil {
		panic("CalendarCommand is nil")
	}
	if calendarQuery == nil {
		panic("CalendarQuery is nil")
	}
	return &CalendarHttpHandler{
		eventCommand:    eventCommand,
		eventQuery:      eventQuery,
		calendarCommand: calendarCommand,
		calendarQuery:   calendarQuery,
	}
}
