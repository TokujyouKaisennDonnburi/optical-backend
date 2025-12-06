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
	GetEventsByDate(ctx context.Context, userId uuid.UUID, date time.Time) ([]output.EventTodayQueryOutputItem, error)
}
