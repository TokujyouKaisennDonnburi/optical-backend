package output

import (
	"database/sql"

	"github.com/google/uuid"
)

// メンバー一覧取得で渡す出力データ
type MembersQueryOutput struct {
	UserId   uuid.UUID
	Name     string
	JoinedAt sql.NullTime
}
