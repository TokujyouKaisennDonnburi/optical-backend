package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/service/command"

type ScheduleHttpHandler struct {
	scheduleCommand *command.ScheduleCommand
}

func NewScheduleHttpHandler(
	createCommandService *command.ScheduleCommand,
) *ScheduleHttpHandler {
	if createCommandService == nil {
		panic("CreateSchedule is nil")
	}
	return &ScheduleHttpHandler{
		scheduleCommand: createCommandService,
	}
}
