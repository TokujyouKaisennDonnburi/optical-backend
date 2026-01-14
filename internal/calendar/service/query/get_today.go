package query

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

type EventTodayQueryInput struct {
	UserId uuid.UUID
}

func (q *EventQuery) GetToday(
	ctx context.Context,
	input EventTodayQueryInput,
) (*output.EventTodayQueryOutput, error) {
	now := time.Now().UTC()
	outputs, err := q.eventRepository.GetEventsByDate(ctx, input.UserId, now)
	if err != nil {
		return nil, err
	}
	return &output.EventTodayQueryOutput{
		Date: now.Format("2006-01-02"),
		Items: outputs,
	}, nil
}
