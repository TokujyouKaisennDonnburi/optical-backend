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

type AnalyzableCalendar struct {
	Id       string             `json:"calendar_id"`
	Name     string             `json:"calendar_name"`
	Color    string             `json:"calendar_color"`
	ImageUrl string             `json:"imageUrl"`
	Members  []AnalyzableMember `json:"members"`
	Options  []AnalyzableOption `json:"options"`
}

type AnalyzableMember struct {
	UserId   string    `json:"user_id"`
	Name     string    `json:"user_name"`
	JoinedAt time.Time `json:"joined_at"`
}

type AnalyzableOption struct {
	Id   int32  `json:"option_id"`
	Name string `json:"option_name"`
}
