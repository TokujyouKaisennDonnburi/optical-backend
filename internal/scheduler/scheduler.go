package scheduler

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type Status int8

const (
	Nil     Status = 0
	Good    Status = 1
	UnKnown Status = 2
	Bad     Status = 3
)

type Scheduler struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Title      string
	Memo       string
	StartTime  time.Time
	EndTime    time.Time
	LimitTime  time.Time
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
	Status       Status
}

func NewScheduler(calendarId, userId uuid.UUID, title, memo string, startTime, endTime, limitTime time.Time, isAllDay bool) (*Scheduler, error) {
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
	s := &Scheduler{
		Id:         id,
		CalendarId: calendarId,
		UserId:     userId,
		Title:      title,
		Memo:       memo,
		StartTime:  startTime,
		EndTime:    endTime,
		LimitTime:  limitTime,
		IsAllDay:   isAllDay,
	}
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
	err = s.SetLimitTime(limitTime)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Scheduler) SetTitle(title string) error {
	titleLength := utf8.RuneCountInString(title)
	if titleLength < 1 || titleLength > 32 {
		return apperr.ValidationError("title is invalid")
	}
	s.Title = title
	return nil
}
func (s *Scheduler) SetMemo(memo string) error {
	memoLength := utf8.RuneCountInString(memo)
	if memoLength > 256 {
		return apperr.ValidationError("memo is invalid")
	}
	s.Memo = memo
	return nil
}
func (s *Scheduler) SetStartEndTime(startTime, endTime time.Time) error {
	if startTime.After(endTime) {
		return apperr.ValidationError("start time before must be end time")
	}
	s.StartTime = startTime
	s.EndTime = endTime
	return nil
}

func (s *Scheduler) SetLimitTime(limitTime time.Time) error {
	if limitTime.IsZero() {
		return apperr.ValidationError("limit time is invalid")
	}
	now := time.Now()
	if limitTime.After(s.StartTime) {
		return apperr.ValidationError("limit time must be before or equal to startTime")
	}
	if limitTime.Before(now) || limitTime.Equal(now) {
		return apperr.ValidationError("limit time must be after current time")
	}
	s.LimitTime = limitTime
	return nil
}

func NewAttendance(id, schedulerId, userId uuid.UUID, comment string) (*SchedulerAttendance, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if schedulerId == uuid.Nil {
		return nil, errors.New("schedulerId is nil")
	}
	if userId == uuid.Nil {
		return nil, errors.New("userId is nil")
	}
	s := &SchedulerAttendance{
		Id:          id,
		SchedulerId: schedulerId,
		UserId:      userId,
		Comment:     comment,
	}
	err = s.SetComment(comment)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SchedulerAttendance) SetComment(comment string) error {
	commentLength := len(comment)
	if commentLength > 255 {
		return apperr.ValidationError("comment length error")
	}
	s.Comment = comment
	return nil
}
func NewStatus(attendanceId uuid.UUID, time time.Time, status int8) (*SchedulerStatus, error) {
	if attendanceId == uuid.Nil {
		return nil, errors.New("attendanceId is nil")
	}
	s := &SchedulerStatus{
		AttendanceId: attendanceId,
		Time:         time,
		Status:       Status(status),
	}
	err := s.SetStatus(status)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *SchedulerStatus) SetStatus(status int8) error {
	if status < 1 || status > 3 {
		return apperr.ValidationError("status is invalid")
	}
	s.Status = Status(status)
	return nil
}
