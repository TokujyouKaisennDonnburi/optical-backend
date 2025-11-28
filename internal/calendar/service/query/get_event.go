package query

import (
	"time"

	"github.com/google/uuid"
)

// イベント取得用の入力データ
type EventQueryInput struct {
	CalendarId uuid.UUID
}

// イベント取得用の出力データ
type EventQueryOutput struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	Title      string
	Memo       string
	Color      string
	AllDay     bool
	StartAt    time.Time
	EndAt      time.Time
	CreatedAt  time.Time
}

// カレンダーに紐づくイベント一覧を取得する
func (q *EventQuery) ListEventsByCalendarId(
	input EventQueryInput,
) ([]EventQueryOutput, error) {
	events, err := q.eventRepository.ListEventsByCalendarId(
		q.ctx,
		input.CalendarId,
	)
	if err != nil {
		return nil, err
	}
	var outputs []EventQueryOutput
	for _, event := range events {
		outputs = append(outputs, EventQueryOutput{
			Id:         event.Id,
			CalendarId: event.CalendarId,
			Title:      event.Title,
			Memo:       event.Memo,
			Color:      event.Color,
			AllDay:     event.ScheduledTime.AllDay,
			StartAt:    event.ScheduledTime.StartTime,
			EndAt:      event.ScheduledTime.EndTime,
			CreatedAt:  event.CreatedAt,
		})
	}
	return outputs, nil
}
