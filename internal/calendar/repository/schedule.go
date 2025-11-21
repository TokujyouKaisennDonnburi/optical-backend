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
		userId uuid.UUID,
		optionIds []uuid.UUID,
		createFn func(*user.User, []option.Option) (*calendar.Calendar, error),
	) error
}
