package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

// Input: Handler層から受け取るデータ
type EventQueryInput struct {
	CalendarID uuid.UUID
	UserID     uuid.UUID
}

// イベント一覧取得
func (q *EventQuery) ListGetEvents(ctx context.Context, input EventQueryInput) ([]output.EventQueryOutput, error) {
	events, err := q.eventRepository.ListEventsByCalendarId(ctx, input.CalendarID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
