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
	memberQuery     *query.MemberQuery
}

func NewCalendarHttpHandler(
	eventCommand *command.EventCommand,
	eventQuery *query.EventQuery,
	calendarCommand *command.CalendarCommand,
	calendarQuery *query.CalendarQuery,
	memberQuery *query.MemberQuery,
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
	if memberQuery == nil {
		panic("MemberQuery is nil")
	}
	return &CalendarHttpHandler{
		eventCommand:    eventCommand,
		eventQuery:      eventQuery,
		calendarCommand: calendarCommand,
		calendarQuery:   calendarQuery,
		memberQuery:     memberQuery,
	}
}
