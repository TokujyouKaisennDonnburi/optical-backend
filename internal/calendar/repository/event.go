package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
)

type EventRepository interface {
	Create(
		ctx context.Context,
		calendarId uuid.UUID,
		createFn func(*calendar.Calendar) (*calendar.Event, error),
	) error

	// 一覧取得
	ListEventsByCalendarId(
		ctx context.Context,
		calendarId uuid.UUID,
	) ([]*calendar.Event, error)
}
