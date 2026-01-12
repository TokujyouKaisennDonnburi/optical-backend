package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (*AgentCommandPsqlRepository) UpdateOptions(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	optionIds []int32,
) error {
	return transact.Transact(ctx, func(tx *sqlx.Tx) error {
		var exists bool
		query := `
			SELECT 1
			FROM calendar_members
			WHERE 
				calendar_members.calendar_id = $1
				AND calendar_members.user_id = $2
		`
		err := tx.GetContext(ctx, &exists, query, calendarId, userId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.NotFoundError("invalid permission")
		}
		query = `
			DELETE FROM calendar_options
			WHERE calendar_options.calendar_id = $1
		`
		_, err = tx.ExecContext(ctx, query, calendarId)
		if err != nil {
			return err
		}
		query = `
			INSERT INTO calendar_options(calendar_id, option_id)
			VALUES (:calendarId, :optionId)
		`
		optionMaps := make([]map[string]any, len(optionIds))
		for i, optionId := range optionIds {
			optionMaps[i] = map[string]any{
				"calendarId": calendarId,
				"optionId":   optionId,
			}
		}
		_, err = tx.NamedExecContext(ctx, query, optionMaps)
		return err
	})
}
