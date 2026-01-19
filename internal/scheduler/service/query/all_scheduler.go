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

func (q *SchedulerQuery) AllScheduler(ctx context.Context, input AllSchedulerInput) ([]scheduler.Scheduler, error) {
	// option check
	options, err := q.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	hasOption := false
	for _, opt := range options {
		if opt.Id == option.OPTION_SCHEDULER {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option not enabled")
	}
	// repository
	results, err := q.schedulerRepository.FindAllSchedulerById(ctx, input.CalendarId, input.UserId)
	if err != nil {
		return nil, err
	}
	return results, nil
}
