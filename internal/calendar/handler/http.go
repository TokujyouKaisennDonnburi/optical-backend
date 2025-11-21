package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"

type CalendarHttpHandler struct {
	calendarCommand *command.CalendarCommand
}

func NewCalendarHttpHandler(
	calendarCommand *command.CalendarCommand,
) *CalendarHttpHandler {
	if calendarCommand == nil {
		panic("CreateSchedule is nil")
	}
	return &CalendarHttpHandler{
		calendarCommand: calendarCommand,
	}
}
