package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
)

type CalendarHttpHandler struct {
	eventCommand    *command.EventCommand
	calendarCommand *command.CalendarCommand
	eventQuery      *query.EventQuery
	calendarQuery   *query.CalendarQuery
}

func NewCalendarHttpHandler(
	eventCommand *command.EventCommand,
	calendarCommand *command.CalendarCommand,
	eventQuery *query.EventQuery,
	calendarQuery *query.CalendarQuery,
) *CalendarHttpHandler {
	if eventCommand == nil {
		panic("EventCommand is nil")
	}
	if calendarCommand == nil {
		panic("CreateSchedule is nil")
	}
	if eventQuery == nil {
		panic("EventQuery is nil")
	}
	if calendarQuery == nil {
		panic("CalendarQuery is nil")
	}
	return &CalendarHttpHandler{
		eventCommand:    eventCommand,
		calendarCommand: calendarCommand,
		eventQuery:      eventQuery,
		calendarQuery:   calendarQuery,
	}
}
