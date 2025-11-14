package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
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
	var err error
	var options []option.Option
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		options, err = FindOptionsByIds(ctx, tx, ids)
		return err
	})
	return options, err
}

func FindOptionsByIds(ctx context.Context, tx *sqlx.Tx, ids []uuid.UUID) ([]option.Option, error) {
	if len(ids) == 0 {
		return []option.Option{}, nil
	}
	query := `
		SELECT id, name
			FROM options
		WHERE 
			id in (?)
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
