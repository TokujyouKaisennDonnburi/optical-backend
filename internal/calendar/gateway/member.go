package gateway


import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MemberPsqlRepository struct {
	db *sqlx.DB
}

func (r *MemberPsqlRepository)Create(ctx context.Context, userId, calendarId uuid.UUID, email string)error{
	// userIdをDBから取得
	query := `
		INSERT INTO calendar_members (calendar_id, user_id)
		SELECT $3, u.id
		FROM users u
		WHERE u.email = $1
		AND EXISTS (
			SELECT 1 
			FROM calendar_members cm 
			WHERE cm.calendar_id = $3 
			AND cm.user_id = $2)
			`
	err := r.db.GetContext(ctx, &userId, query, email, userId, calendarId)
	if err != nil {
		return err
	}
	// memberをDBに挿入
	query =`
		INSERT INTO calendar_members(calendar_id,member_id)
		VALUES ($1, $2)
	`
	_, err = r.db.ExecContext(ctx, query, calendarId, ???)
	if err != nil {
		return err
	}
	// member
	return nil
}

