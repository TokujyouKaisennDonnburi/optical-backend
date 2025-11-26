package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type CalendarRepository interface {
	Create(
		ctx context.Context,
		userId, imageId uuid.UUID,
		optionIds []uuid.UUID,
		createFn func(*user.User, *calendar.Image, []option.Option) (*calendar.Calendar, error),
	) error
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]calendar.Calendar, error)
}
