package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
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

func (r *OptionPsqlRepository) FindByIds(ctx context.Context, ids []int32) ([]option.Option, error) {
	var err error
	var options []option.Option
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		options, err = psql.FindOptionsByIds(ctx, tx, ids)
		return err
	})
	return options, err
}

