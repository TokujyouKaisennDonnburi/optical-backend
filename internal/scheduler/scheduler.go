package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	Title string
	Memo string
	StartTime time.Time
	EndTitle time.Time
	IsAllDay bool
}

