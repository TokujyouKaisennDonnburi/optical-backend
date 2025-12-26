package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OptionPsqlRepository struct {
	db *sqlx.DB
}

func NewOptionPsqlRepository(db *sqlx.DB) *OptionPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &OptionPsqlRepository{
		db: db,
	}
}

func (r *OptionPsqlRepository) FindAll(ctx context.Context) ([]option.Option, error) {
	var models []psql.OptionModel
	query := `
		SELECT id, name FROM options 
		WHERE deprecated = FALSE
	`
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, err
	}
	options := make([]option.Option, len(models))
	for i, model := range models {
		options[i] = option.Option{
			Id:         model.Id,
			Name:       model.Name,
			Deprecated: false,
		}
	}
	return options, err
}

func (r *OptionPsqlRepository) FindOptionsByIds(ctx context.Context, optionIds []int32) ([]option.Option, error) {
	if len(optionIds) == 0 {
		return []option.Option{}, nil
	}
	query := `
		SELECT id, name
		FROM options
		WHERE id = ANY($1)
		AND deprecated = FALSE
		ORDER BY id
	`
	var models []psql.OptionModel
	err := r.db.SelectContext(ctx, &models, query, pq.Array(optionIds))
	if err != nil {
		return nil, err
	}
	options := make([]option.Option, len(models))
	for i, model := range models {
		options[i] = option.Option{
			Id:   model.Id,
			Name: model.Name,
		}
	}
	return options, nil
}
