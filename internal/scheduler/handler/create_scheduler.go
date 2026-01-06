package handler

import (
	"time"

	"github.com/google/uuid"
)

type SchedulerCreateRequest struct {
	Title     string    `json:"title"`
	Memo      string    `json:"memo"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	LimitTime time.Time `json:"limitTime"`
	IsAllDay  bool      `json:"isAllDay"`
}
type SchedulerCreateResponse struct {
	Id         uuid.UUID `json:"id"`
	CalendarId uuid.UUID`json:"calendarId"`
	Title      string    `json:"title"`
	Memo       string    `json:"memo"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	IsAllDay   bool      `json:"isAllDay"`
}

func SchedulerCreate
