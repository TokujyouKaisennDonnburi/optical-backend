package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerResultInput struct {
	SchedulerId uuid.UUID
	UserId      uuid.UUID
}

func (q *SchedulerQuery) SchedulerResult(ctx context.Context, input SchedulerResultInput) (*output.SchedulerResultOutput, error) {
	// repository
	result, err := q.schedulerRepository.FindResultById(ctx, input.SchedulerId, input.UserId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
