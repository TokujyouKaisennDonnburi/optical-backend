package output

import (
	"time"

	"github.com/google/uuid"
)

// イベント取得で渡す出力データ
type EventQueryOutput struct {
	Id            uuid.UUID
	CalendarId    uuid.UUID
	UserId        uuid.UUID
	CalendarColor string
	Title         string
	Memo          string
	Location      string
	IsAllDay      bool
	StartAt       string
	EndAt         string
	CreatedAt     string
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
	Location      string
	Memo          string
	StartAt       time.Time
	EndAt         time.Time
	IsAllDay      bool
}

type EventMonthQueryOutput struct {
	Date  string
	Items []EventTodayQueryOutputItem
}

type EventMonthQueryOutputItem struct {
	CalendarId    uuid.UUID
	CalendarName  string
	CalendarColor string
	Id            uuid.UUID
	Title         string
	Location      string
	Memo          string
	StartAt       time.Time
	EndAt         time.Time
	IsAllDay      bool
}

// 検索結果用の出力データ
type EventSearchQueryOutput struct {
	Items []EventSearchQueryOutputItem // 検索結果の配列
	Total int                          // マッチ全件数
	Limit int                          // 1ページ当たりの表示件数
}

type EventSearchQueryOutputItem struct {
	CalendarId    uuid.UUID
	CalendarName  string
	CalendarColor string
	EventId       uuid.UUID
	EventTitle    string
	Location      string
	Memo          string
	StartAt       time.Time
	EndAt         time.Time
	IsAllDay      bool
}
