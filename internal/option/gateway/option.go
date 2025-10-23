package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

type OptionModel struct {
	id   uuid.UUID
	name string
}

func (r *OptionPsqlRepository) FindByIds(ctx context.Context, ids []uuid.UUID) ([]option.Option, error) {
	query := `
		SELECT id, name
			FROM options
		WHERE 
			id in (?)
	`
	optionModels := []OptionModel{}
	err := r.db.SelectContext(ctx, optionModels, query, ids)
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
