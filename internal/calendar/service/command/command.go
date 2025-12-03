package command

import (
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
)

type CalendarCommand struct {
	calendarRepository calendarRepo.CalendarRepository
	optionRepository   optionRepo.OptionRepository
	imageRepository    calendarRepo.ImageRepository
	memberRepository calendarRepo.MemberRepository
}

type EventCommand struct {
	eventRepository calendarRepo.EventRepository
}

func NewCalendarCommand(
	calendarRepository calendarRepo.CalendarRepository,
	optionRepository optionRepo.OptionRepository,
	imageRepository calendarRepo.ImageRepository,
	memberRepository calendarRepo.MemberRepository,
) *CalendarCommand {
	if calendarRepository == nil {
		panic("calendarRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	if imageRepository == nil {
		panic("imageRepository is nil")
	}
	if memberRepository == nil {
		panic("memberRepository is nil")
	}
	return &CalendarCommand{
		calendarRepository: calendarRepository,
		optionRepository:   optionRepository,
		imageRepository:    imageRepository,
		memberRepository:   memberRepository,

	}
}

func NewEventCommand(eventRepository calendarRepo.EventRepository) *EventCommand {
	if eventRepository == nil {
		panic("eventRepository is nil")
	}
	return &EventCommand{
		eventRepository: eventRepository,
	}
}

