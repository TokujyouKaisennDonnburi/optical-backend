package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/service/command"

type ScheduleHttpHandler struct {
	createCommandService *command.CreateSchedule
}

func NewScheduleHttpHandler(
	createCommandService *command.CreateSchedule,
) *ScheduleHttpHandler {
	if createCommandService == nil {
		panic("CreateSchedule is nil")
	}
	return &ScheduleHttpHandler{
		createCommandService: createCommandService,
	}
}
