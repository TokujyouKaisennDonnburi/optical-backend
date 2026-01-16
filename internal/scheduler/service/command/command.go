package command

import (
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	schedulerRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/repository"
)

type SchedulerCommand struct {
	schedulerRepository schedulerRepo.SchedulerRepository
	optionRepository    optionRepo.OptionRepository
	eventRepository     calendarRepo.EventRepository
}

func NewSchedulerCommand(
	schedulerRepository schedulerRepo.SchedulerRepository,
	optionRepository optionRepo.OptionRepository,
	eventRepository calendarRepo.EventRepository,
) *SchedulerCommand {
	if schedulerRepository == nil {
		panic("schedulerRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	if eventRepository == nil {
		panic("eventRepository is nil")
	}
	return &SchedulerCommand{
		schedulerRepository: schedulerRepository,
		optionRepository:    optionRepository,
		eventRepository:     eventRepository,
	}
}
