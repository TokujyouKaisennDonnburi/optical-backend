package agent

import (
	"time"
)

type AnalyzableEvent struct {
	CalendarId    string    `json:"calendar_id"`
	CalendarName  string    `json:"calendar_name"`
	CalendarColor string    `json:"calendar_color"`
	Id            string    `json:"event_id"`
	Title         string    `json:"event_title"`
	Location      string    `json:"location"`
	Memo          string    `json:"memo"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
	IsAllday      bool      `json:"is_allday"`
}
