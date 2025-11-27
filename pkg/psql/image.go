package psql

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ImageModel struct {
	id  uuid.UUID
	url string
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
		return nil, err
	}
	return &calendar.Image{
		Id: imageModel.id,
		Url: imageModel.url,
		Valid: true,
	}, nil
}
