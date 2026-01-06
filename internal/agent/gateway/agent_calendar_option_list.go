package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OptionModel struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
}

func (r *AgentQueryPsqlRepository) FindOptionsByCalendarId(
	ctx context.Context,
	userId, calendarId uuid.UUID,
) ([]agent.AnalyzableOption, error) {
	var options []agent.AnalyzableOption
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
			SELECT options.id, options.name
			FROM options
			WHERE id IN 
			(
				SELECT option_id
				FROM calendar_options
				JOIN calendar_members
					ON calendar_members.calendar_id = calendar_options.calendar_id
				WHERE 
					calendar_options.calendar_id = $1
					AND calendar_members.user_id = $2
			)
		`
		var models []OptionModel
		err := tx.SelectContext(ctx, &models, query, userId, calendarId)
		if err != nil {
			return err
		}
		for _, model := range models {
			options = append(options, agent.AnalyzableOption{
				Id:   model.Id,
				Name: model.Name,
			})
		}
		return nil
	})
	return options, err
}
