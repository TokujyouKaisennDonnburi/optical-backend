package gateway

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MemberPsqlRepository struct {
	db *sqlx.DB
}

func (r *MemberPsqlRepository)Create(ctx context.Context, calendarId uuid.UUID, email string)(*calendar.Member, error){
	// TODO: 1. emailからusersテーブルでユーザー情報を取得する
	query := `
		SELECT u.id
		FROM users u
		WHERE email = $1
	`
	// TODO: 2. ユーザーが見つからなかった場合はエラーを返す
	if query = nil {
		return nil, errors.New("user not found")
	}
	// TODO: 3. calendar_membersテーブルに挿入する
	query :=`
		INSERT INTO calendar_members(calendar_id,member_id)
		VALUES (calendarId,userId)
	`
	// TODO: 4. 取得したユーザー情報からMemberを作成して返す
	calendar.NewMember(userId uuid.UUID, name string)(*calendar.Member, error)
}

