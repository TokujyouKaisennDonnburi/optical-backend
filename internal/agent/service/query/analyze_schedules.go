package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type AnalyzeSchedulesQueryInput struct {
	UserInput   string
	UserId      uuid.UUID
	StreamingFn func(context.Context, []byte) error
}

func (q *AgentQuery) AnalyzeSchedules(ctx context.Context, input AnalyzeSchedulesQueryInput) error {
	if input.UserInput == "" {
		return apperr.ValidationError("invalid user prompt")
	}
	schedules, err := q.eventRepository.FindAnalyzableEventsByUserId(ctx, input.UserId)
	if err != nil {
		return err
	}
	err = q.agentRepository.AnalyzeSchedules(
		ctx,
		input.UserInput,
		schedules,
		input.StreamingFn,
	)
	return err
}
