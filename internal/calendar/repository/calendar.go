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
		optionIds []int32,
		createFn func(*calendar.Image, []calendar.Member, []option.Option) (*calendar.Calendar, error),
	) error

	Update(
		ctx context.Context,
		imageId uuid.UUID,
		memberEmails []string,
		optionIds []int32,
		updateFn func(*calendar.Image, []calendar.Member, []option.Option) (*calendar.Calendar, error),
	) error

	FindByUserId(ctx context.Context, userId uuid.UUID) ([]output.CalendarListQueryOutput, error)
	FindByUserCalendarId(ctx context.Context, userId, calendarId uuid.UUID) (*calendar.Calendar, error)
}
