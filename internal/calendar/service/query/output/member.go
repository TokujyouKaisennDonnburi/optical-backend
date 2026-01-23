package output

import (
	"time"

	"github.com/google/uuid"
)

// 参加メンバー一覧取得で渡す出力データ
type ParticipantsMembersQueryOutput struct {
	UserId   uuid.UUID
	Name     string
	JoinedAt time.Time
}
