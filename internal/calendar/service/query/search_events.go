package query

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

// handler層から受け取るデータ
type SearchEventQueryInput struct {
	UserId    uuid.UUID
	Query     string
	StartFrom time.Time
	StartTo   time.Time
	Limit     int
	Offset    int
}

// イベント検索
func (q *EventQuery) SearchEvents(
	ctx context.Context,
	input SearchEventQueryInput,
) (*output.EventSearchQueryOutput, error) {
	params := repository.SearchEventsParams{
		UserId:    input.UserId,
		Query:     input.Query,
		StartFrom: input.StartFrom,
		StartTo:   input.StartTo,
		Limit:     input.Limit,
		Offset:    input.Offset,
	}

	return q.eventRepository.SearchEvents(ctx, params)
}
