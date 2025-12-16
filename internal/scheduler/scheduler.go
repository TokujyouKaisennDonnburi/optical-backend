package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Title      string
	Memo       string
	StartTime  time.Time
	EndTitle   time.Time
	IsAllDay   bool
}

type Scheduler_attendance struct {
	Id        uuid.UUID
	CalendaId uuid.UUID
	UserId    uuid.UUID
	Coment    string
}
type Scheduler_status struct {
	AttendanceId uuid.UUID
	Time         time.Time
	Status       int
}


