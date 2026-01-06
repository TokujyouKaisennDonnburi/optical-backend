package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type SchedulerRepository interface {
	CreateScheduler(ctx context.Context, id, calendcarId uuid.UUID, title, memo string, startTime, endTime, limitTime time.Time, isAllDay bool)(scheduler.Scheduler, error)
}
