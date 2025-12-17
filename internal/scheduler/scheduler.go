package scheduler

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Scheduler struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Title      string
	Memo       string
	StartTime  time.Time
	EndTime    time.Time
	IsAllDay   bool
}

type Scheduler_attendance struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Comment    string
}

type Scheduler_status struct {
	AttendanceId uuid.UUID
	Time         time.Time
	Status       int
}

func NewScheduler(userId, calendarId uuid.UUID, title, memo string, startTime, endTime time.Time, isAllDay bool) (*Scheduler, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if calendarId == uuid.Nil {
		return nil, errors.New("calendarId is nil")
	}
	if userId == uuid.Nil {
		return nil, errors.New("userId is nil")
	}
	s, err := &Scheduler{
		Id:         id,
		CalendarId: calendarId,
		UserId:     userId,
		Title:      title,
		Memo:       memo,
		StartTime:  startTime,
		EndTime:    endTime,
		IsAllDay:   isAllDay,
	}, nil
	err = s.SetTitle(title)
	if err != nil {
		return nil, err
	}
	err = s.SetMemo(memo)
	if err != nil {
		return nil, err
	}
	err = s.SetStartEndTime(startTime, endTime)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Scheduler) SetTitle(title string) error {
	titleLength := utf8.RuneCountInString(title)
	if titleLength < 1 || titleLength > 32 {
		return errors.New("title is invalid")
	}
	s.Title = title
	return nil
}
func (s *Scheduler) SetMemo(memo string) error {
	memoLength := utf8.RuneCountInString(memo)
	if memoLength < 1 || memoLength > 256 {
		return errors.New("memo is invalid")
	}
	s.Memo = memo
	return nil
}
func (s *Scheduler) SetStartEndTime(startTime, endTime time.Time) error {
	if startTime.After(endTime){
		return errors.New("")
	}
	return nil
}
