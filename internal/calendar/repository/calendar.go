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
		calendar *calendar.Calendar,
	) error

	Update(
		ctx context.Context,
		userId uuid.UUID,
		calendarId uuid.UUID,
		optionIds []int32,
		updateFn func(
			existingCalendar *calendar.Calendar,
			options []option.Option,
		) (*calendar.Calendar, error),
	) error
	Delete(
		ctx context.Context,
		calendarId uuid.UUID,
		userId uuid.UUID,
	) error

	FindById(ctx context.Context, id uuid.UUID) (*calendar.Calendar, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]output.CalendarListQueryOutput, error)
	FindByUserCalendarId(ctx context.Context, userId, calendarId uuid.UUID) (*calendar.Calendar, error)
}
