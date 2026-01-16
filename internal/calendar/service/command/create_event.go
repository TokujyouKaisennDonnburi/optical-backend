package command

import (
	"context"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
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
	event, err := calendar.NewEvent(input.CalendarId, input.Title, input.Memo, input.Location, *scheduledTime)
	if err != nil {
		return nil, err
	}
	// リポジトリに保存
	err = c.eventRepository.Create(ctx, event.CalendarId, func(cal *calendar.Calendar) (*calendar.Event, error) {
		// ユーザーがカレンダーのメンバーかチェック
		for _, member := range cal.Members {
			if member.UserId == input.UserId {
				return event, nil
			}
		}
		return nil, errors.New("User is not member.")
	})
	if err != nil {
		return nil, err
	}
	return &EventCreateOutput{
		Id: event.Id,
	}, nil
}
