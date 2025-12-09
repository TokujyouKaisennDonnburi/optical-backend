package query

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type EventThisMonthQueryInput struct {
	UserId uuid.UUID
}

func (q *EventQuery) GetThisMonth(
	ctx context.Context,
	input EventThisMonthQueryInput,
) (*output.EventMonthQueryOutput, error) {
	now := time.Now()
	outputs, err := q.eventRepository.GetEventsByMonth(ctx, input.UserId, now)
	if err != nil {
		return nil, err
	}
	return &output.EventMonthQueryOutput{
		Date:  now.Format("2006-01"),
		Items: outputs,
	}, nil
}

type EventMonthQueryInput struct {
	UserId uuid.UUID
	Month  string
}

func (q *EventQuery) GetByMonth(
	ctx context.Context,
	input EventMonthQueryInput,
) (*output.EventMonthQueryOutput, error) {
	month, err := time.Parse("2006-01", input.Month)
	if err != nil {
		return nil, apperr.ValidationError("invalid month")
	}
	outputs, err := q.eventRepository.GetEventsByMonth(ctx, input.UserId, month)
	if err != nil {
		return nil, err
	}
	return &output.EventMonthQueryOutput{
		Date:  input.Month,
		Items: outputs,
	}, nil
}
