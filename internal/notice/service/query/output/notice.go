package output

import (
	"github.com/google/uuid"
)

// 通知取得で渡す出力データ
type NoticeQueryOutput struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	EventId    *uuid.UUID
	CalendarId *uuid.UUID
	Title      string
	Content    string
	IsRead     bool
	CreatedAt  string
}
