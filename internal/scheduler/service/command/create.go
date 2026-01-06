package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type SchedulerCreateInput struct {
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Title      string
	Memo       string
	StartTime  time.Time
	EndTime    time.Time
	LimitTime  time.Time
	IsAllDay   bool
}

type SchedulerCreateOutput struct {
	Id         uuid.UUID
}

func (c *SchedulerCommand) CreateScheduler(ctx context.Context, input SchedulerCreateInput) (*scheduler.Scheduler, error) {
	scheduler, err := scheduler.NewScheduler(ctx, input.CalendarId, input.UserId, input.Title, input.Memo, input.StartTime, input.EndTime, input.LimitTime, input.IsAllDay)
	if err != nil {
		return nil, err
	}
	result, err := c.schedulerRepository.CreateScheduler(ctx, scheduler.Id, scheduler.CalendarId, scheduler.Title, scheduler.Memo, scheduler.StartTime, scheduler.EndTime, scheduler.LimitTime, scheduler.IsAllDay)
	if err != nil {
		return nil, err
	}
	return &SchedulerCreateOutput{
		Id: result.Id,
	}, nil
}
