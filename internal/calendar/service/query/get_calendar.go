package query

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)

type GetCalendarInput struct {
	CalendarId uuid.UUID
	UserId     uuid.UUID
}

type CalendarQueryOutput struct {
	Id      uuid.UUID
	Name    string
	Color   string
	Image   calendar.Image
	Members []calendar.Member
	Options []option.Option
}

// カレンダー情報を取得
func (q *CalendarQuery) GetCalendar(ctx context.Context, input GetCalendarInput)(*CalendarQueryOutput, error) {
	calendar, err := q.calendarRepository.FindById(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	isMember := false
	for _, m := range calendar.Members {
		if m.UserId == input.UserId{
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.New("user is not in members")
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
