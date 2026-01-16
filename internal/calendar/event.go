package calendar

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MIN_EVENT_TITLE_LENGTH    = 1
	MAX_EVENT_TITLE_LENGTH    = 32
	MAX_EVENT_LOCATION_LENGTH = 32
	MAX_EVENT_MEMO_LENGTH     = 255
)

type Event struct {
	Id            uuid.UUID
	CalendarId    uuid.UUID
	Title         string
	Memo          string
	Location      string
	ScheduledTime ScheduledTime
}

type ScheduledTime struct {
	AllDay    bool
	StartTime time.Time
	EndTime   time.Time
}

func NewEvent(calendarId uuid.UUID, title, memo, location string, scheduledTime ScheduledTime) (*Event, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("Event `id` is nil")
	}
	if calendarId == uuid.Nil {
		return nil, errors.New("Evnet `calendarId` is nil")
	}
	event := Event{
		Id:         id,
		CalendarId: calendarId,
	}
	err = event.SetTitle(title)
	if err != nil {
		return nil, err
	}
	err = event.SetMemo(memo)
	if err != nil {
		return nil, err
	}
	err = event.SetLocation(location)
	if err != nil {
		return nil, err
	}
	err = event.SetScheduledTime(scheduledTime)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *Event) SetTitle(title string) error {
	titleLength := utf8.RuneCountInString(title)
	if titleLength < MIN_EVENT_TITLE_LENGTH || titleLength > MAX_EVENT_TITLE_LENGTH {
		return errors.New("Event `title` length is invalid")
	}
	e.Title = title
	return nil
}

func (e *Event) SetMemo(memo string) error {
	memoLength := utf8.RuneCountInString(memo)
	if memoLength > MAX_EVENT_MEMO_LENGTH {
		return errors.New("Event `memo` length is invalid")
	}
	e.Memo = memo
	return nil
}

func (e *Event) SetLocation(location string) error {
	locationLength := utf8.RuneCountInString(location)
	if locationLength > MAX_EVENT_LOCATION_LENGTH {
		return errors.New("Event `location` length is invalid")
	}
	e.Location = location
	return nil
}

func (e *Event) SetScheduledTime(scheduledTime ScheduledTime) error {
	if scheduledTime.IsZero() {
		return errors.New("Event `scheduledTime` is zero")
	}
	e.ScheduledTime = scheduledTime
	return nil
}

func NewScheduledTime(allDay bool, startTime, endTime time.Time) (*ScheduledTime, error) {
	if allDay {
		return &ScheduledTime{
			StartTime: startTime,
			EndTime: endTime,
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
