package gateway

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

// トークンで招待を取得
func (r *MemberPsqlRepository) FindInvitationByToken(ctx context.Context, token uuid.UUID) (*calendar.Invitation, error) {
	query := `
		SELECT id, calendar_id, email, joined_user_id, token, expires_at, used_at, created_at
		FROM calendar_invitations
		WHERE token = $1
	`
	var row invitationRow
	err := r.db.GetContext(ctx, &row, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFoundError("invitation not found")
		}
		return nil, err
	}
	return row.toInvitation(), nil
}

// カレンダーIDで招待一覧を取得(未使用かつ有効期限内)
func (r *MemberPsqlRepository) FindPendingInvitationsByCalendarId(ctx context.Context, calendarId uuid.UUID) ([]*calendar.Invitation, error) {
	query := `
		SELECT id, calendar_id, email, joined_user_id, token, expires_at, used_at, created_at
		FROM calendar_invitations
		WHERE calendar_id = $1
		AND used_at IS NULL
		AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	var rows []invitationRow
	err := r.db.SelectContext(ctx, &rows, query, calendarId)
	if err != nil {
		return nil, err
	}
	invitations := make([]*calendar.Invitation, len(rows))
	for i, row := range rows {
		invitations[i] = row.toInvitation()
	}
	return invitations, nil
}
