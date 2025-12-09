package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type CalendarRepository interface {
	Create(
		ctx context.Context,
		imageId uuid.UUID,
		memberEmails []string,
		optionIds []uuid.UUID,
		createFn func(*calendar.Image, []calendar.Member, []option.Option) (*calendar.Calendar, error),
	) error
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]output.CalendarQueryOutput, error)
}
