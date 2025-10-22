package command

import (
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	scheduleRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/repository"
)

type ScheduleCommand struct {
	scheduleRepository scheduleRepo.ScheduleRepository
	optionRepository   optionRepo.OptionRepository
}

