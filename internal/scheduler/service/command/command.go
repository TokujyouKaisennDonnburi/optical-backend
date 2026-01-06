package command

import "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/repository"

type SchedulerCommand struct {
	schedulerRepository repository.SchedulerRepository
}

func NewSchedulerCommand(schedulerRepository repository.SchedulerRepository) *SchedulerCommand {
	if schedulerRepository == nil {
		panic("schedulerRepository is nil")
	}
	return &SchedulerCommand{
		schedulerRepository: schedulerRepository,
	}
}
