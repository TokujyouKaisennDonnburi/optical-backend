package command

import (
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/repository"
)

type SchedulerCommand struct {
	schedulerRepository repository.SchedulerRepository
	optionRepository    optionRepo.OptionRepository
}

func NewSchedulerCommand(schedulerRepository repository.SchedulerRepository, optionRepository optionRepo.OptionRepository) *SchedulerCommand {
	if schedulerRepository == nil {
		panic("schedulerRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	return &SchedulerCommand{
		schedulerRepository: schedulerRepository,
		optionRepository:    optionRepository,
	}
}
