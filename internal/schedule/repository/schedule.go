package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Create(
		ctx context.Context,
		userId uuid.UUID,
		optionIds []uuid.UUID,
		createFn func(*user.User, []option.Option) (*schedule.Schedule, error),
	) error
}
