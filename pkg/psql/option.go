package psql

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/jmoiron/sqlx"
)

type OptionModel struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
}

func FindOptionsByIds(ctx context.Context, tx *sqlx.Tx, ids []int32) ([]option.Option, error) {
	if len(ids) == 0 {
		return []option.Option{}, nil
	}
	query := `
		SELECT id, name
			FROM options
		WHERE 
			id in (?)
		AND deprecated = FALSE
		ORDER BY id
	`
	optionModels := []OptionModel{}
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	query = tx.Rebind(query)
	err = tx.SelectContext(ctx, &optionModels, query, args...)
	if err != nil {
		return nil, err
	}
	options := []option.Option{}
	for _, optionModel := range optionModels {
		options = append(options, option.Option{
			Id:   optionModel.Id,
			Name: optionModel.Name,
		})
	}
	return options, nil
}
