package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerQueryInput struct {
	SchedulerId uuid.UUID
	UserId      uuid.UUID
	CalendarId  uuid.UUID
}

func (q *SchedulerQuery) SchedulerQuery(
	ctx context.Context,
	input SchedulerQueryInput,
) (*output.SchedulerOutput, error) {
	// repository
	result, err := q.schedulerRepository.FindSchedulerById(ctx, input.SchedulerId)
	if err != nil {
		return nil, err
	}
	// assign
	return &output.SchedulerOutput{
		Id:           result.Id,
		CalendarId:   result.CalendarId,
		UserId:       result.UserId,
		Title:        result.Title,
		Memo:         result.Memo,
		LimitTime:    result.LimitTime,
		IsAllDay:     result.IsAllDay,
		PossibleDate: result.PossibleDate,
	}, nil
}
