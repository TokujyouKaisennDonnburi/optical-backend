package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"

type CalendarHttpHandler struct {
	eventCommand    *command.EventCommand
	calendarCommand *command.CalendarCommand
}

func NewCalendarHttpHandler(
	eventCommand *command.EventCommand,
	calendarCommand *command.CalendarCommand,
) *CalendarHttpHandler {
	if eventCommand == nil {
		panic("EventCommand is nil")
	}
	if calendarCommand == nil {
		panic("CreateSchedule is nil")
	}
	return &CalendarHttpHandler{
		eventCommand:    eventCommand,
		calendarCommand: calendarCommand,
	}
}
