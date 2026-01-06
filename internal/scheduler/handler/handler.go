package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/command"

type SchedulerHttpHandler struct {
	schedulerCommand *command.SchedulerCommand
}

func NewSchedulerHttpHandler(schedulerCommand *command.SchedulerCommand) *SchedulerHttpHandler {
	if schedulerCommand == nil {
		panic("SchedulerCommand is nil")
	}
	return &SchedulerHttpHandler{
		schedulerCommand: schedulerCommand,
	}
}
