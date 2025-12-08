package gateway

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MemberPsqlRepository struct {
	db *sqlx.DB
}

func NewMemberPsqlRepository(db *sqlx.DB) *MemberPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &MemberPsqlRepository{
		db: db,
	}
}

func (r *MemberPsqlRepository) Create(ctx context.Context, userId, calendarId uuid.UUID, emails []string) error {
	query := `
		INSERT INTO calendar_members (calendar_id, user_id, joined_at)
		SELECT $3, u.id, NULL
		FROM users u
		WHERE u.email = ANY($1)
		AND EXISTS (
			SELECT 1 FROM calendar_members cm 
			WHERE cm.calendar_id = $3 
			AND cm.user_id = $2
		)
		AND NOT EXISTS (
			SELECT 1 FROM calendar_members cm
			WHERE cm.calendar_id = $3
			AND cm.user_id = u.id
		)
		`
	result, err := r.db.ExecContext(ctx, query, pq.Array(emails), userId, calendarId)
	if err != nil {
		return err
	}
	// 実行できている行数をとってくる
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("already member or not member for me")
	}
	return nil
}

