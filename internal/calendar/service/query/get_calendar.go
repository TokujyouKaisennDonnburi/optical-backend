package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
	"github.com/google/uuid"
)

type GetCalendarInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type CalendarQueryOutput struct {
	Id       uuid.UUID
	Name     string
	Color    calendar.Color
	Image    calendar.Image
	ImageUrl string
	Members  []calendar.Member
	Options  []option.Option
}

// カレンダー情報を取得
func (q *CalendarQuery) GetCalendar(ctx context.Context, input GetCalendarInput) (*CalendarQueryOutput, error) {
	calendar, err := q.calendarRepository.FindByUserCalendarId(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return nil, err
	}
	imageUrl := ""
	if calendar.Image.Valid {
		imageUrl = storage.GetImageStorageBaseUrl() + "/" + calendar.Image.Url
	}
	return &CalendarQueryOutput{
		Id:       calendar.Id,
		Name:     calendar.Name,
		Color:    calendar.Color,
		Image:    calendar.Image,
		ImageUrl: imageUrl,
		Members:  calendar.Members,
		Options:  calendar.Options,
	}, nil
}
