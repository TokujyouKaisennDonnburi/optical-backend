package gateway

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
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

func (r *MemberPsqlRepository) Create(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
	emailModel := make([]string, len(emails))
	for i, e := range emails {
		emailModel[i] = string(e)
	}
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
	result, err := r.db.ExecContext(ctx, query, pq.Array(emailModel), userId, calendarId)
	if err != nil {
		return err
	}
	// 実行できている行数
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("already member or not member for me")
	}
	return nil
}

func (r *MemberPsqlRepository) Join(ctx context.Context, userId, calendarId uuid.UUID) error {
	query := `
	UPDATE calendar_members
	SET joined_at = NOW()
	WHERE user_id = $1
	AND calendar_id = $2
	AND joined_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, userId, calendarId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("not invited or already joined")
	}
	return nil
}

