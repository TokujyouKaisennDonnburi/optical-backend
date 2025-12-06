package output

import (
	"time"

	"github.com/google/uuid"
)

type EventTodayQueryOutput struct {
	Date  string
	Items []EventTodayQueryOutputItem
}

type EventTodayQueryOutputItem struct {
	CalendarId    uuid.UUID
	CalendarName  string
	CalendarColor string
	Id            uuid.UUID
	Title         string
	Color         string
	Location      string
	Memo          string
	StartAt       time.Time
	EndAt         time.Time
	IsAllDay      bool
}
