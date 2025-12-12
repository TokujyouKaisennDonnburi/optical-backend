package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
)

type CalendarListQueryInput struct {
	UserId uuid.UUID
}

// ユーザーが所属するカレンダー一覧を取得する
func (q *CalendarQuery) GetCalendars(ctx context.Context, input CalendarListQueryInput) ([]output.CalendarListQueryOutput, error) {
	return q.calendarRepository.FindByUserId(ctx, input.UserId)
}
