package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/notice"
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/transact"
)

type CalendarCommand struct {
	transactor            transact.TransactionProvider
	calendarRepository    calendarRepo.CalendarRepository
	optionRepository      optionRepo.OptionRepository
	imageRepository       calendarRepo.ImageRepository
	memberRepository      calendarRepo.MemberRepository
	emailRepository       calendarRepo.EmailRepository
	calendarNoticeService *notice.CalendarNoticeService
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
	calendarNoticeService *notice.CalendarNoticeService,
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
	if calendarNoticeService == nil {
		panic("calendarNoticeService is nil")
	}
	return &CalendarCommand{
		transactor:            transactor,
		calendarRepository:    calendarRepository,
		optionRepository:      optionRepository,
		imageRepository:       imageRepository,
		memberRepository:      memberRepository,
		emailRepository:       emailRepository,
		calendarNoticeService: calendarNoticeService,
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
