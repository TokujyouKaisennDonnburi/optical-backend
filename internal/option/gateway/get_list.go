package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OptionListQueryModel struct {
	OptionId   int32  `db:"id"`
	Name       string `db:"name"`
	Deprecated bool   `db:"deprecated"`
}

func (r *OptionPsqlRepository) GetList(ctx context.Context, userId uuid.UUID) ([]option.Option, error) {
	query := `
	SELECT options.id, options.name, options.deprecated
	FROM options
	WHERE EXISTS (SELECT 1 FROM users WHERE users.id = $1)
	ORDER BY options.id
	`
	var rows []OptionListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, err
	}
	result := make([]option.Option, len(rows))
	for i, row := range rows {
		result[i] = option.Option{
			Id:         row.OptionId,
			Name:       row.Name,
			Deprecated: row.Deprecated,
		}
	}
	return result, nil
}

func (r *OptionPsqlRepository) FindsByIds(ctx context.Context, ids []int32) ([]option.Option, error) {
	var rows []OptionListQueryModel
	err := db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		query := `
			SELECT id, name
				FROM options
			WHERE 
				id = ANY ($1)
			AND deprecated = FALSE
			ORDER BY id
		`
		err := tx.SelectContext(ctx, &rows, query, pq.Array(ids))
		return err
	})
	if err != nil {
		return nil, err
	}
	result := make([]option.Option, len(rows))
	for i, row := range rows {
		result[i] = option.Option{
			Id:         row.OptionId,
			Name:       row.Name,
			Deprecated: false,
		}
	}
	return result, nil
}
