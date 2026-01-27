package command

import (
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/transact"
)

type CalendarCommand struct {
	transactor         transact.TransactionProvider
	calendarRepository calendarRepo.CalendarRepository
	optionRepository   optionRepo.OptionRepository
	imageRepository    calendarRepo.ImageRepository
	memberRepository   calendarRepo.MemberRepository
	emailRepository    calendarRepo.EmailRepository
}

type EventCommand struct {
	transactor         transact.TransactionProvider
	eventRepository    calendarRepo.EventRepository
	calendarRepository calendarRepo.CalendarRepository
}

func NewCalendarCommand(
	transactor transact.TransactionProvider,
	calendarRepository calendarRepo.CalendarRepository,
	optionRepository optionRepo.OptionRepository,
	imageRepository calendarRepo.ImageRepository,
	memberRepository calendarRepo.MemberRepository,
	emailRepository calendarRepo.EmailRepository,
) *CalendarCommand {
	if transactor == nil {
		panic("transactor is nil")
	}
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
	if emailRepository == nil {
		panic("emailRepository is nil")
	}
	return &CalendarCommand{
		transactor:         transactor,
		calendarRepository: calendarRepository,
		optionRepository:   optionRepository,
		imageRepository:    imageRepository,
		memberRepository:   memberRepository,
		emailRepository:    emailRepository,
	}
}

func NewEventCommand(
	transactor transact.TransactionProvider,
	eventRepository calendarRepo.EventRepository,
	calendarRepository calendarRepo.CalendarRepository,
) *EventCommand {
	if transactor == nil {
		panic("transactor is nil")
	}
	if eventRepository == nil {
		panic("eventRepository is nil")
	}
	if calendarRepository == nil {
		panic("calendarRepository is nil")
	}
	return &EventCommand{
		transactor: transactor,
		eventRepository: eventRepository,
		calendarRepository: calendarRepository,
	}
}
