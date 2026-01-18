package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type AllSchedulerInput struct {
	CalendarId uuid.UUID
	UserId     uuid.UUID
}

func (q *SchedulerQuery) AllScheduler(ctx context.Context, input SchedulerQueryInput) (*scheduler.Scheduler, error) {
	// option check
	options, err := q.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
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
	// repository
	result, err := q.schedulerRepository.FindAllSchedulerById(ctx, input.CalendarId, input.UserId)
	if err != nil {
		return nil, err
	}
	// assign
	return &scheduler.Scheduler{
		Id:         result.Id,
		CalendarId: result.CalendarId,
		UserId:     result.UserId,
		Title:      result.Title,
		Memo:       result.Memo,
		LimitTime:  result.LimitTime,
		IsAllDay:   result.IsAllDay,
	}, nil
}
