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

type OptionListQueryModel struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	Deprecated bool   `db:"deprecated"`
}

func (r *OptionPsqlRepository) GetList(ctx context.Context, userId uuid.UUID) ([]*option.Option, error) {
	query := `
	SELECT options.id, options.name, options.deprecated
	FROM options
	LEFT JOIN users
	ON users.id = $1
	`
	var rows []OptionListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, err
	}
	result := make([]option.Option, len(rows))
	for i, row := range rows {
		result[i] = option.Option{
			Id:         row.Id,
			Name:       row.Name,
			Deprecated: row.Deprecated,
		}
	}
	return result, nil
}
