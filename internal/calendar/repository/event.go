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
	Update(
		ctx context.Context,
		eventId uuid.UUID,
		updateFn func(*calendar.Event) (*calendar.Event, error),
	) error
}
