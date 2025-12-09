package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
)

type EventUpdateInput struct {
	UserId    uuid.UUID
	EventId   uuid.UUID
	Title     string
	Memo      string
	Color     string
	Location  string
	StartTime time.Time
	EndTime   time.Time
	IsAllDay  bool
}

// 予定を更新する
func (c *EventCommand) UpdateEvent(ctx context.Context, input EventUpdateInput) error {
	// 予定時間を作成
	scheduledTime, err := calendar.NewScheduledTime(input.IsAllDay, input.StartTime, input.EndTime)
	if err != nil {
		return err
	}
	err = c.eventRepository.Update(ctx, input.UserId, input.EventId, func(event *calendar.Event) (*calendar.Event, error) {
		err = event.SetTitle(input.Title)
		if err != nil {
			return nil, err
		}
		err = event.SetColor(input.Color)
		if err != nil {
			return nil, err
		}
		err = event.SetMemo(input.Memo)
		if err != nil {
			return nil, err
		}
		err = event.SetLocation(input.Location)
		if err != nil {
			return nil, err
		}
		err = event.SetScheduledTime(*scheduledTime)
		if err != nil {
			return nil, err
		}
		return event, nil
	})
	return err
}
