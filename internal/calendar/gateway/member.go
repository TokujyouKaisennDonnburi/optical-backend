package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MemberPsqlRepository struct {
	db *sqlx.DB
}

func (r *MemberPsqlRepository)Create(ctx context.Context, calendarId uuid.UUID, email string)(*calendar.Calendar, error){
	// TODO: 返り値の型をcalendar.Memberに修正する（今はcalendar.Calendarになってる）

	// TODO: 1. emailからusersテーブルでユーザー情報を取得する
	//   SELECT id, name FROM users WHERE email = $1

	// TODO: 2. ユーザーが見つからなかった場合はエラーを返す

	// TODO: 3. calendar_membersテーブルに挿入する
	//   INSERT INTO calendar_members(calendar_id, user_id) VALUES ($1, $2)

	// TODO: 4. 取得したユーザー情報からMemberを作成して返す
	//   calendar.NewMember(userId, name) を使う

	// TODO: memberを作る
	query := `
	SELECT u.id
	FROM users u
	INSERT INTO calendar_members(calendar_id,user_id)
	VALUES (:calendarId, :u.id)
	`
}

