package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
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
	Id uuid.UUID
}

func (c *SchedulerCommand) CreateScheduler(ctx context.Context, input SchedulerCreateInput) (*SchedulerCreateOutput, error) {
	// option check
	options, err := c.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	hasOption := false
	for _, x := range options {
		if x.Id == option.OPTION_SCHEDULER {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option not enabled")
	}
	// domain
	scheduler, err := scheduler.NewScheduler(input.CalendarId, input.UserId, input.Title, input.Memo, input.StartTime, input.EndTime, input.LimitTime, input.IsAllDay)
	if err != nil {
		return nil, err
	}
	// repository
	result, err := c.schedulerRepository.CreateScheduler(ctx, scheduler.Id, scheduler.CalendarId, scheduler.UserId, scheduler.Title, scheduler.Memo, scheduler.StartTime, scheduler.EndTime, scheduler.LimitTime, scheduler.IsAllDay)
	if err != nil {
		return nil, err
	}
	// output
	return &SchedulerCreateOutput{
		Id: result.Id,
	}, nil
}
