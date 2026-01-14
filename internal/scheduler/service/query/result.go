package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerResultInput struct {
	CalendarId  uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
}

func (q *SchedulerQuery) SchedulerResult(ctx context.Context, input SchedulerResultInput) (*output.SchedulerResultOutput, error) {
	// repository
	result, err := q.schedulerRepository.FindByMemberId(ctx, input.CalendarId, input.SchedulerId, input.UserId)
	if err != nil {
		return nil, err
	}
	// assign
	return &output.SchedulerResultOutput{
		UserId:    result.UserId,
		Title:     result.Title,
		Memo:      result.Memo,
		LimitTime: result.LimitTime,
		IsAllDay:  result.IsAllDay,
		Members:   result.Members,
		Date:      result.Date,
	}, nil
}
