package agent

import (
	"encoding/json"
	"time"
)

type AnalyzableEvent struct {
	CalendarId    string
	CalendarName  string
	CalendarColor string
	Id            string
	Title         string
	Location      string
	Memo          string
	StartAt       time.Time
	EndAt         time.Time
	IsAllDay      bool
}

type analyzableEventJson struct {
	CalendarId    string `json:"calendarId"`
	CalendarName  string `json:"calendarName"`
	CalendarColor string `json:"calendarColor"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Location      string `json:"location"`
	Memo          string `json:"memo"`
	StartAt       string `json:"startAt"`
	EndAt         string `json:"endAt"`
	IsAllDay      bool   `json:"isAllDay"`
}

func (e AnalyzableEvent) MarshalJSON() ([]byte, error) {
	jsonBody := analyzableEventJson{
		CalendarId:    e.CalendarId,
		CalendarName:  e.CalendarName,
		CalendarColor: e.CalendarColor,
		Id:            e.Id,
		Title:         e.Title,
		Location:      e.Location,
		Memo:          e.Memo,
		StartAt:       e.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:         e.EndAt.Format("2006-01-02 15:04:05"),
		IsAllDay:      e.IsAllDay,
	}
	return json.Marshal(jsonBody)
}
