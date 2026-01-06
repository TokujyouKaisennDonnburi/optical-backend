package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
)

type AgentQueryRepository interface {
	FindEventByUserIdAndDate(
		ctx context.Context,
		userId uuid.UUID,
		startAt, endAt time.Time,
	) ([]agent.AnalyzableEvent, error)
	FindCalendarEventByUserIdAndDate(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		startAt, endAt time.Time,
	) ([]agent.AnalyzableEvent, error)
	FindCalendarByUserId(
		ctx context.Context,
		userId uuid.UUID,
	) ([]agent.AnalyzableCalendar, error)
	FindCalendarByIdAndUserId(
		ctx context.Context,
		userId, calendarId uuid.UUID,
	) (*agent.AnalyzableCalendar, error)
	FindOptionsByCalendarId(
		ctx context.Context,
		userId, calendarId uuid.UUID,
	) ([]agent.AnalyzableOption, error)
}

type AgentCommandRepository interface {
	CreateEvents(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		createFn func(*calendar.Calendar) ([]calendar.Event, error),
	) error
	UpdateOptions(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		optionIds []int32,
	) error
}
