package output

import (
	"time"

	"github.com/google/uuid"
)

type SchedulerAttendanceQuery struct {
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
