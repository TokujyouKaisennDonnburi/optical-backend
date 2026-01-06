package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type SchedulerCreateInput struct {
	CalendarId uuid.UUID
	UserId uuid.UUID
	Title     string
	Memo      string
	StartTime time.Time
	EndTime   time.Time
	LimitTime time.Time
	IsAllDay  bool
}

func (c *SchedulerCommand) CreateScheduler(ctx context.Context, input SchedulerCreateInput)(*scheduler.Scheduler, error){
	scheduler, err := scheduler.NewScheduler(ctx, input.CalendarId, input.UserId, input.Title,input.Memo, input.StartTime, input.EndTime, input.LimitTime, input.IsAllDay)
	if err != nil {
		return nil, err
	}
}

