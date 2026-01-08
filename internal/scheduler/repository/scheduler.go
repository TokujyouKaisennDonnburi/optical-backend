package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type SchedulerRepository interface {
	CreateScheduler(ctx context.Context, id, calendarId, userId uuid.UUID, title, memo string, possibleDates []scheduler.PossibleDate, limitTime time.Time, isAllDay bool)error
}
