package query

import (
	scheduler "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/repository"
	option "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
)

type SchedulerQuery struct {
	schedulerRepository scheduler.SchedulerRepository
	optionRepository    option.OptionRepository
}

func NewSchedulerQuery(schedulerRepository scheduler.SchedulerRepository, optionRepository option.OptionRepository) *SchedulerQuery {
	if schedulerRepository == nil {
		panic("schedulerRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	return &SchedulerQuery{
		schedulerRepository: schedulerRepository,
		optionRepository:    optionRepository,
	}
}
