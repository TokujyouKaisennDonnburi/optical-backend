package query

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

var ErrCalendarNotBelongToUser = errors.New("calendar does not belong to user")

// Input: Handler層から受け取るデータ
type EventQueryInput struct {
	CalendarID uuid.UUID
	UserID     uuid.UUID
}

// イベント一覧取得
func (q *EventQuery) ListGetEvents(ctx context.Context, input EventQueryInput) ([]output.EventQueryOutput, error) {
	// 権限チェック: カレンダーがユーザーに属しているか確認
	exists, err := q.eventRepository.ExistsCalendarByUserIdAndCalendarId(ctx, input.UserID, input.CalendarID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrCalendarNotBelongToUser
	}

	events, err := q.eventRepository.ListEventsByCalendarId(ctx, input.CalendarID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
