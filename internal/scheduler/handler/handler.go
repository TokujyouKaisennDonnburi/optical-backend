package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query"
)

type SchedulerHttpHandler struct {
	schedulerCommand *command.SchedulerCommand
	schedulerQuery   *query.SchedulerQuery
}

func NewSchedulerHttpHandler(schedulerCommand *command.SchedulerCommand, schedulerQuery *query.SchedulerQuery) *SchedulerHttpHandler {
	if schedulerCommand == nil {
		panic("SchedulerCommand is nil")
	}
	if schedulerQuery == nil {
		panic("SchedulerQuery is nil")
	}
	return &SchedulerHttpHandler{
		schedulerCommand: schedulerCommand,
		schedulerQuery:   schedulerQuery,
	}
}
