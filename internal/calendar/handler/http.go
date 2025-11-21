package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"

type CalendarHttpHandler struct {
	calendarCommand *command.ScheduleCommand
}

func NewCalendarHttpHandler(
	calendarCommand *command.ScheduleCommand,
) *CalendarHttpHandler {
	if calendarCommand == nil {
		panic("CreateSchedule is nil")
	}
	return &CalendarHttpHandler{
		calendarCommand: calendarCommand,
	}
}
