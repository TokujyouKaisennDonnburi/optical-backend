package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
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
	) ([]output.EventQueryOutput, error)

	// カレンダーが指定されたユーザーに属しているかチェック
	ExistsCalendarByUserIdAndCalendarId(
		ctx context.Context,
		userId uuid.UUID,
		calendarId uuid.UUID,
	) (bool, error)
	GetEventsByDate(ctx context.Context, userId uuid.UUID, date time.Time) ([]output.EventTodayQueryOutputItem, error)
}
