package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
	"github.com/google/uuid"
)

type CalendarListQueryInput struct {
	UserId uuid.UUID
}

// ユーザーが所属するカレンダー一覧を取得する
func (q *CalendarQuery) GetCalendars(ctx context.Context, input CalendarListQueryInput) ([]output.CalendarListQueryOutput, error) {
	calendars, err := q.calendarRepository.FindByUserId(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	for i, cal := range calendars {
		if cal.Image.Valid {
			calendars[i].ImageUrl = storage.GetImageStorageBaseUrl() + "/" + cal.Image.Url
		}
	}
	return calendars, nil
}
