package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/jmoiron/sqlx"
)

type OptionModel struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
}

func (r *AgentQueryPsqlRepository) FindOptions(
	ctx context.Context,
) ([]agent.AnalyzableOption, error) {
	var options []agent.AnalyzableOption
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
			SELECT options.id, options.name
			FROM options
			WHERE options.deprecated = FALSE 
		`
		var models []OptionModel
		err := tx.SelectContext(ctx, &models, query)
		if err != nil {
			return err
		}
		for _, model := range models {
			option := agent.AnalyzableOption{
				Id:   model.Id,
				Name: model.Name,
			}
			if desc, ok := getDescriptionMap()[option.Id]; ok {
				option.Description = desc
			}
			options = append(options, option)
		}
		return nil
	})
	return options, err
}
