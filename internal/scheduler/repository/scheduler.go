package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerRepository interface {
	CreateScheduler(ctx context.Context, id, calendarId, userId uuid.UUID, title, memo string, possibleDates []scheduler.PossibleDate, limitTime time.Time, isAllDay bool) error
	FindSchedulerById(ctx context.Context, id uuid.UUID) (*output.SchedulerOutput, error)
	FindAllSchedulerById(ctx context.Context, id uuid.UUID) (*scheduler.Scheduler, error)
	AddAttendance(ctx context.Context, id, schedulerId, userId uuid.UUID, comment string, schedulerStatus []scheduler.SchedulerStatus) error
	FindResultByIdAndUserId(ctx context.Context, schedulerId, userId uuid.UUID) (*output.SchedulerResultOutput, error)
	FindAttendanceById(ctx context.Context, calendarId, schedulerId, userId uuid.UUID) (*output.SchedulerAttendanceOutput, error)
}
