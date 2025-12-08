package output

import (
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
