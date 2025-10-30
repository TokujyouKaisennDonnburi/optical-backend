package command

import (
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	scheduleRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/repository"
)

type ScheduleCommand struct {
	scheduleRepository scheduleRepo.ScheduleRepository
	optionRepository   optionRepo.OptionRepository
}

func NewScheduleCommand(scheduleRepository scheduleRepo.ScheduleRepository, optionRepository optionRepo.OptionRepository) *ScheduleCommand {
	if scheduleRepository == nil {
		panic("scheduleRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	return &ScheduleCommand{
		scheduleRepository: scheduleRepository,
		optionRepository: optionRepository,
	}
}

