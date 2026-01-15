package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerUserStatusInput struct {
	CalendarId  uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
}

func (q *SchedulerQuery) UserStatusQuery(ctx context.Context, input SchedulerUserStatusInput) (*output.SchedulerUserOutput, error) {
	result, err := q.schedulerRepository.FindStatusById(ctx, input.CalendarId, input.SchedulerId, input.UserId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
