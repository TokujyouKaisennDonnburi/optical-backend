package psql

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/jmoiron/sqlx"
)

type OptionModel struct {
	id   int32
	name string
}

func FindOptionsByIds(ctx context.Context, tx *sqlx.Tx, ids []int) ([]option.Option, error) {
	if len(ids) == 0 {
		return []option.Option{}, nil
	}
	query := `
		SELECT id, name
			FROM options
		WHERE 
			id in (?)
		AND deprecated = FALSE
	`
	optionModels := []OptionModel{}
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	err = tx.SelectContext(ctx, optionModels, query, args...)
	if err != nil {
		return nil, err
	}
	options := []option.Option{}
	for _, optionModel := range optionModels {
		options = append(options, option.Option{
			Id:   optionModel.id,
			Name: optionModel.name,
		})
	}
	return options, nil
}

