package gateway

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MemberPsqlRepository struct {
	db *sqlx.DB
}

type MemberQueryModel struct {
	UserId   uuid.UUID    `db:"user_id"`
	UserName string       `db:"user_name"`
	JoinedAt sql.NullTime `db:"joined_at"`
}

func NewMemberPsqlRepository(db *sqlx.DB) *MemberPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &MemberPsqlRepository{
		db: db,
	}
}

func (r *MemberPsqlRepository) Invite(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
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
		return apperr.ForbiddenError("already member or not member for me")
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
		return apperr.NotFoundError("not invited or already joined")
	}
	return nil
}

func (r *MemberPsqlRepository) Reject(ctx context.Context, userId, calendarId uuid.UUID) error {
	query := `
	DELETE FROM calendar_members
	WHERE calendar_members.user_id = $1
	AND calendar_members.calendar_id = $2
	AND calendar_members.joined_at IS NULL
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
		return apperr.NotFoundError("user or calendar or joined not found")
	}
	return nil
}

// カレンダーメンバー一覧取得
func (r *MemberPsqlRepository) FindMembers(ctx context.Context, calendarId uuid.UUID) ([]output.MembersQueryOutput, error) {
	query := `
		SELECT
			cm.user_id,
			u.name AS user_name,
			cm.joined_at
		FROM calendar_members cm
		INNER JOIN users u ON u.id = cm.user_id
		WHERE cm.calendar_id = $1
		ORDER BY cm.joined_at ASC
	`
	var rows []MemberQueryModel
	err := r.db.SelectContext(ctx, &rows, query, calendarId)
	if err != nil {
		return nil, err
	}

	members := make([]output.MembersQueryOutput, len(rows))
	for i, row := range rows {
		members[i] = output.MembersQueryOutput{
			UserId:   row.UserId,
			Name:     row.UserName,
			JoinedAt: row.JoinedAt,
		}
	}
	return members, nil
}
