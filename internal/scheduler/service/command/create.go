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
	LimitTime  time.Time
	IsAllDay   bool
	Dates      []SchedulerCreateDateInput
}

type SchedulerCreateDateInput struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
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
	if len(input.Dates) == 0 {
		return nil, apperr.ValidationError("dates is empty")
	}
	// domain
	schedulerEntity, err := scheduler.NewScheduler(input.CalendarId, input.UserId, input.Title, input.Memo, input.LimitTime, input.IsAllDay)
	if err != nil {
		return nil, err
	}
	possibleDates := make([]scheduler.PossibleDate, len(input.Dates))
	for i, date := range input.Dates {
		possibleDate, err := scheduler.NewPossibleDate(date.Date, date.StartTime, date.EndTime, input.IsAllDay)
		if err != nil {
			return nil, err
		}
		possibleDates[i] = *possibleDate
	}
	// repository
	result, err := c.schedulerRepository.CreateScheduler(
		ctx,
		schedulerEntity.Id,
		schedulerEntity.CalendarId,
		schedulerEntity.UserId,
		schedulerEntity.Title,
		schedulerEntity.Memo,
		possibleDates,
		schedulerEntity.LimitTime,
		schedulerEntity.IsAllDay,
	)
	if err != nil {
		return nil, err
	}
	// output
	return &SchedulerCreateOutput{
		Id: result.Id,
	}, nil
}
