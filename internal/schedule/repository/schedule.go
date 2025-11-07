package repository

import(
	"context"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
)

  type ScheduleRepository interface {
      FindAll(ctx context.Context) ([]*schedule.Schedule, error)
  }

