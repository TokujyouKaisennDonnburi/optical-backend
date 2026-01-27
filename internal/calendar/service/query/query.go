package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
)

type EventQuery struct {
	eventRepository repository.EventRepository
}

func NewEventQuery(eventRepo repository.EventRepository) *EventQuery {
	return &EventQuery{
		eventRepository: eventRepo,
	}
}

type CalendarQuery struct {
	calendarRepository repository.CalendarRepository
}

func NewCalendarQuery(calendarRepository repository.CalendarRepository) *CalendarQuery {
	if calendarRepository == nil {
		panic("CalendarRepository is nil")
	}
	return &CalendarQuery{
		calendarRepository: calendarRepository,
	}
}

// メンバー関連のクエリサービス
type MemberQuery struct {
	memberRepository repository.MemberRepository
}

// MemberQueryのコンストラクタ
func NewMemberQuery(memberRepository repository.MemberRepository) *MemberQuery {
	if memberRepository == nil {
		panic("MemberRepository is nil")
	}
	return &MemberQuery{
		memberRepository: memberRepository,
	}
}
