package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/repository"

type SchedulerQuery struct {
	schedulerRepository repository.SchedulerRepository
}

func NewSchedulerQuery(schedulerRepository repository.SchedulerRepository) *SchedulerQuery {
	if schedulerRepository == nil {
		panic("schedulerRepository is nil")
	}
	return &SchedulerQuery{
		schedulerRepository: schedulerRepository,
	}
}

