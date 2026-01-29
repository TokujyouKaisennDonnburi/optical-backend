package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type EventCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Title      string
	Memo       string
	Location   string
	StartTime  time.Time
	EndTime    time.Time
	IsAllDay   bool
}

type EventCreateOutput struct {
	Id uuid.UUID
}

// 予定を新規作成する
func (c *EventCommand) Create(ctx context.Context, input EventCreateInput) (*EventCreateOutput, error) {
	// 予定時間を作成
	scheduledTime, err := calendar.NewScheduledTime(input.IsAllDay, input.StartTime, input.EndTime)
	if err != nil {
		return nil, err
	}
	// 予定を作成
	event, err := calendar.NewEvent(input.CalendarId, input.UserId, input.Title, input.Memo, input.Location, *scheduledTime)
	if err != nil {
		return nil, err
	}
	err = c.transactor.Transact(ctx, func(ctx context.Context) error {
		calendar, err := c.calendarRepository.FindByUserCalendarId(ctx, input.UserId, input.CalendarId)
		if err != nil {
			return err
		}
		exists := false
		for _, member := range calendar.Members {
			if member.UserId == input.UserId {
				exists = true
				break
			}
		}
		if !exists {
			return apperr.ForbiddenError("user is not in calendar members")
		}
		// リポジトリに保存
		err = c.eventRepository.Create(ctx, calendar.Id, event)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &EventCreateOutput{
		Id: event.Id,
	}, nil
}
