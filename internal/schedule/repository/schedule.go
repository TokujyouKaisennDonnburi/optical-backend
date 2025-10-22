package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *schedule.Schedule) error
}
