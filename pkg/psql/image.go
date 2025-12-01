package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ImageModel struct {
	Id  uuid.UUID `db:"id"`
	Url string    `db:"url"`
}

func FindImageById(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*calendar.Image, error) {
	query := `
		SELECT id, url
		FROM calendar_images
		WHERE id = $1
	`
	var imageModel ImageModel
	err := tx.GetContext(ctx, &imageModel, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &calendar.Image{
				Valid: false,
			}, nil
		}
		return nil, err
	}
	return &calendar.Image{
		Id:    imageModel.Id,
		Url:   imageModel.Url,
		Valid: true,
	}, nil
}
