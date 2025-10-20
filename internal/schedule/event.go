package schedule

import (
	"encoding/hex"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MIN_EVENT_TITLE_LENGTH = 1
	MAX_EVENT_TITLE_LENGTH = 32
	MAX_EVENT_MEMO_LENGTH  = 255
)

type Event struct {
	Id            uuid.UUID
	Title         string
	Memo          string
	Color         string
	ScheduledTime ScheduledTime
}

type ScheduledTime struct {
	AllDay    bool
	StartTime time.Time
	EndTime   time.Time
}

func NewEvent(title, memo, color string, scheduledTime ScheduledTime) (*Event, error) {
	// id
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("Event `id` is nil")
	}
	// title
	titleLength := utf8.RuneCountInString(title)
	if titleLength < MIN_EVENT_TITLE_LENGTH || titleLength > MAX_EVENT_TITLE_LENGTH {
		return nil, errors.New("Event `title` length is invalid")
	}
	// memo
	memoLength := utf8.RuneCountInString(memo)
	if memoLength > MAX_EVENT_MEMO_LENGTH {
		return nil, errors.New("Event `memo` length is invalid")
	}
	// color
	colorLen := utf8.RuneCountInString(color)
	if colorLen != 6 {
		return nil, errors.New("Color length is invalid")
	}
	_, err = hex.DecodeString(color)
	if err != nil {
		return nil, errors.New("Color format is invalid")
	}
	// scheduledTime
	if scheduledTime.IsZero() {
		return nil, errors.New("Event `scheduledTime` is zero")
	}
	return &Event{
		Id:            id,
		Title:         title,
		Color:         color,
		Memo:          memo,
		ScheduledTime: scheduledTime,
	}, nil
}

func NewScheduledTime(allDay bool, startTime, endTime time.Time) (*ScheduledTime, error) {
	if allDay {
		return &ScheduledTime{
			AllDay: true,
		}, nil
	}
	if startTime.IsZero() {
		return nil, errors.New("ScheduledTime `startTime` is zero")
	}
	if endTime.IsZero() {
		return nil, errors.New("ScheduledTime `endTime` is zero")
	}
	return &ScheduledTime{
		AllDay:    false,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

func (s ScheduledTime) IsZero() bool {
	return s == ScheduledTime{}
}
