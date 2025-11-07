package handler

import(
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/service/command"
)

type ScheduleHttpHandler struct {
	scheduleCommand *command.ScheduleCommand
}

func NewScheduleHttpHandler(scheduleCommand *command.ScheduleCommand) *ScheduleHttpHandler {
	if scheduleCommand == nil{
		panic("ScheduleCommand is nil")
	}
	return &ScheduleHttpHandler{
		scheduleCommand: scheduleCommand,
	}
}
