package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type AttendanceQueryInput struct {
	CalendarId  uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
}

func (q *SchedulerQuery) AttendanceQuery(ctx context.Context, input AttendanceQueryInput) (*output.SchedulerAttendanceOutput, error) {
	result, err := q.schedulerRepository.FindAttendanceById(ctx, input.CalendarId, input.SchedulerId, input.UserId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
