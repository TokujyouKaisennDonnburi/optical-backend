package query

import (
	"context"
	"time"

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
	now := time.Now().UTC()
	startFrom := input.StartFrom
	startTo := input.StartTo
	limit := input.Limit
	offset := input.Offset

	// 検索範囲のデフォルト
	if startFrom.IsZero() {
		startFrom = now.AddDate(-3, 0, 0) // 3年前
	}
	if startTo.IsZero() {
		startTo = now.AddDate(3, 0, 0) // 3年後
	}

	// ページネーションのデフォルト
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return q.eventRepository.SearchEvents(
		ctx,
		input.UserId,
		input.Query,
		startFrom,
		startTo,
		limit,
		offset,
	)
}
