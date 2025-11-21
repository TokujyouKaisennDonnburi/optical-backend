package command

import (
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
)

type CalendarCommand struct {
	calendarRepository calendarRepo.CalendarRepository
	optionRepository   optionRepo.OptionRepository
}

func NewCalendarCommand(calendarRepository calendarRepo.CalendarRepository, optionRepository optionRepo.OptionRepository) *CalendarCommand {
	if calendarRepository == nil {
		panic("scheduleRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	return &CalendarCommand{
		calendarRepository: calendarRepository,
		optionRepository: optionRepository,
	}
}

