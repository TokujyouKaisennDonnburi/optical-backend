package output

import (
	"time"

	"github.com/google/uuid"
)

type SchedulerOutput struct {
	Id           uuid.UUID
	CalendarId   uuid.UUID
	UserId       uuid.UUID
	Title        string
	Memo         string
	LimitTime    time.Time
	IsAllDay     bool
	PossibleDate []PossibleDateOutput
}

type PossibleDateOutput struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}

type SchedulerResultOutput struct {
	OwnerId   uuid.UUID
	Title     string
	Memo      string
	LimitTime time.Time
	IsAllDay  bool
	Members   []MemberOutput
	Date      []DateOutput
}

type MemberOutput struct {
	UserId   uuid.UUID
	UserName string
}

type DateOutput struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}
type SchedulerAttendanceOutput struct {
	UserId  uuid.UUID
	Comment string
	Status  []StatusOutput
}
type StatusOutput struct {
	Date   time.Time
	Status int8
}
