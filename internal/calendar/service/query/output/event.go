package output

import (
	"time"

	"github.com/google/uuid"
)

// イベント取得で渡す出力データ
type EventQueryOutput struct {
	Id         uuid.UUID
	CalendarId uuid.UUID
	Title      string
	Memo       string
	Color      string
	Location   string
	IsAllDay   bool
	StartAt    string
	EndAt      string
	CreatedAt  string
}

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
