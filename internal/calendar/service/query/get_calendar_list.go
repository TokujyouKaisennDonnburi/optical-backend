package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
)

type CalendarQueryInput struct {
	UserId uuid.UUID
}



// ユーザーが所属するカレンダー一覧を取得する
func (q *CalendarQuery) GetCalendars(ctx context.Context, input CalendarQueryInput) ([]output.CalendarQueryOutput, error) {
	calendars, err := q.calendarRepository.FindByUserId(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	items := make([]output.CalendarQueryOutput, len(calendars))
	for i, cal := range calendars {
		items[i] = output.CalendarQueryOutput{
			Id:    cal.Id,
			Name:  cal.Name,
			Color: cal.Color,
		}
	}

	return items, nil
}
