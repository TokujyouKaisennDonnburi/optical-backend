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
	// デフォルト値の適用
	startFrom, startTo, limit, offset := applySearchDefaults(
		input.StartFrom,
		input.StartTo,
		input.Limit,
		input.Offset,
	)

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

// 検索パラメータのデフォルト値を適用
func applySearchDefaults(startFrom, startTo time.Time, limit, offset int) (time.Time, time.Time, int, int) {
	now := time.Now()

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

	return startFrom, startTo, limit, offset
}
