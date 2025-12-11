package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
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
