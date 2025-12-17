package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type GetCalendarInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type CalendarQueryOutput struct {
	Id      uuid.UUID
	Name    string
	Color   calendar.Color
	Image   calendar.Image
	Members []calendar.Member
	Options []option.Option
}

// カレンダー情報を取得
func (q *CalendarQuery) GetCalendar(ctx context.Context, input GetCalendarInput)(*CalendarQueryOutput, error) {
	calendar, err := q.calendarRepository.FindByUserCalendarId(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return nil, err
	}
	return &CalendarQueryOutput{
		Id:      calendar.Id,
		Name:    calendar.Name,
		Color:   calendar.Color,
		Image:   calendar.Image,
		Members: calendar.Members,
		Options: calendar.Options,
	}, nil
} 
