package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
)

type CalendarHttpHandler struct {
	eventCommand    *command.EventCommand
	calendarCommand *command.CalendarCommand
	eventQuery      *query.EventQuery
}

func NewCalendarHttpHandler(
	eventCommand *command.EventCommand,
	calendarCommand *command.CalendarCommand,
	eventQuery *query.EventQuery,
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
	return &CalendarHttpHandler{
		eventCommand:    eventCommand,
		calendarCommand: calendarCommand,
		eventQuery:      eventQuery,
	}
}
