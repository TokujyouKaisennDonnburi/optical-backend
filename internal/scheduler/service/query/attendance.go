package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type AttendanceQueryInput struct {
	SchedulerId uuid.UUID
	UserId      uuid.UUID
	CalendarId  uuid.UUID
}

func (q *SchedulerQuery) AttendanceQuery(
	ctx context.Context,
	input AttendanceQueryInput,
) (*output.SchedulerAttendanceQuery, error) {
	// repository
	result, err := q.schedulerRepository.FindById(ctx, input.SchedulerId)
	if err != nil {
		return nil, err
	}
	// assign
	return &output.SchedulerAttendanceQuery{
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
