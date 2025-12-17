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

type SchedulerAttendance struct {
	Id          uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
	Comment     string
}

type SchedulerStatus struct {
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
	if memoLength > 256 {
		return errors.New("memo is invalid")
	}
	s.Memo = memo
	return nil
}
func (s *Scheduler) SetStartEndTime(startTime, endTime time.Time) error {
	if startTime.After(endTime) {
		return errors.New("")
	}
	return nil
}

func NewAttendance(id, schedulerId, userId uuid.UUID, comment string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}
	if schedulerId == uuid.Nil {
		return errors.New("schedulerId is nil")
	}
	if userId == uuid.Nil {
		return errors.New("userId is nil")
	}
	s := &SchedulerAttendance{
		Id:          id,
		SchedulerId: schedulerId,
		UserId:      userId,
		Comment:     comment,
	}
	err = s.SetComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *SchedulerAttendance) SetComment(comment string) error {
	commentLength := len(comment)
	if commentLength > 255 {
		return nil
	}
	return nil
}
func NewStatus(attendanceId uuid.UUID, time time.Time, status int) error {
	if attendanceId == uuid.Nil {
		return errors.New("attendanceId is nil")
	}
	s := &SchedulerStatus{
		AttendanceId: attendanceId,
		Time:         time,
		Status:       status,
	}
	err := s.SetStatus(status)
	if err != nil {
		return err
	}
	return nil
}
func (s *SchedulerStatus) SetStatus(status int) error {
	if status < 0 || status > 2 {
		return errors.New("status is nil")
	}
	return nil
}
