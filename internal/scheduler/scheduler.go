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
	Id            uuid.UUID
	CalendarId    uuid.UUID
	UserId        uuid.UUID
	Title         string
	Memo          string
	LimitTime     time.Time
	IsAllDay      bool
}
type PossibleDate struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
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

func NewScheduler(calendarId, userId uuid.UUID, title, memo string, limitTime time.Time, isAllDay bool) (*Scheduler, error) {
	// id
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// calendarId
	if calendarId == uuid.Nil {
		return nil, errors.New("calendarId is nil")
	}
	// userId
	if userId == uuid.Nil {
		return nil, errors.New("userId is nil")
	}
	s := &Scheduler{
		Id:         id,
		CalendarId: calendarId,
		UserId:     userId,
		Title:      title,
		Memo:       memo,
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
func (s *Scheduler) SetLimitTime(limitTime time.Time) error {
	if limitTime.IsZero() {
		return apperr.ValidationError("limit time is invalid")
	}
	now := time.Now()
	if limitTime.Before(now) || limitTime.Equal(now) {
		return apperr.ValidationError("limit time must be after current time")
	}
	s.LimitTime = limitTime
	return nil
}

func NewPossibleDate(date, startTime, endTime time.Time, isAllDay bool) (*PossibleDate, error) {
	if date.IsZero() {
		return nil, apperr.ValidationError("date is invalid")
	}
	if isAllDay {
		if !isStartOfDay(startTime) {
			return nil, apperr.ValidationError("start time must be start of day")
		}
		if !isEndOfDay(endTime) {
			return nil, apperr.ValidationError("end time must be end of day")
		}
		return &PossibleDate{
			Date:      date,
			StartTime: startTime,
			EndTime:   endTime,
		}, nil
	}
	if startTime.IsZero() {
		return nil, apperr.ValidationError("start time is invalid")
	}
	if endTime.IsZero() {
		return nil, apperr.ValidationError("end time is invalid")
	}
	if !startTime.Before(endTime) {
		return nil, apperr.ValidationError("start time must be before end time")
	}
	return &PossibleDate{
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

func isStartOfDay(t time.Time) bool {
	h, m, s := t.Clock()
	return h == 0 && m == 0 && s == 0 && t.Nanosecond() == 0
}
func isEndOfDay(t time.Time) bool {
	h, m, s := t.Clock()
	return h == 23 && m == 59 && s == 59 && t.Nanosecond() == 0
}

func NewAttendance(schedulerId, userId uuid.UUID, comment string) (*SchedulerAttendance, error) {
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
	a := &SchedulerAttendance{
		Id:          id,
		SchedulerId: schedulerId,
		UserId:      userId,
		Comment:     comment,
	}
	err = a.SetComment(comment)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *SchedulerAttendance) SetComment(comment string) error {
	commentLength := len(comment)
	if commentLength > 255 {
		return apperr.ValidationError("comment length error")
	}
	a.Comment = comment
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
