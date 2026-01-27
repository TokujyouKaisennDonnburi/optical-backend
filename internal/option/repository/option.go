package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type OptionRepository interface {
	FindAll(ctx context.Context) ([]option.Option, error)
	GetList(ctx context.Context, userId uuid.UUID) ([]option.Option, error)
	FindsByIds(ctx context.Context, ids []int32) ([]option.Option, error)
	FindsByCalendarId(ctx context.Context, calendarId uuid.UUID) ([]option.Option, error)
}
