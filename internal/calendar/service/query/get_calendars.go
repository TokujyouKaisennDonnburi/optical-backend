package query

import (
	"context"

	"github.com/google/uuid"
)

type CalendarQueryInput struct {
	UserId uuid.UUID
}

type CalendarItem struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}

type CalendarListOutput struct {
	Calendars []CalendarItem
}

// ユーザーが所属するカレンダー一覧を取得する
func (q *CalendarQuery) GetCalendars(ctx context.Context, input CalendarQueryInput) (*CalendarListOutput, error) {
	calendars, err := q.calendarRepository.FindByUserId(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	items := make([]CalendarItem, len(calendars))
	for i, cal := range calendars {
		items[i] = CalendarItem{
			Id:    cal.Id,
			Name:  cal.Name,
			Color: cal.Color,
		}
	}

	return &CalendarListOutput{
		Calendars: items,
	}, nil
}
