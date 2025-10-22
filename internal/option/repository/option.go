package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type OptionRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]option.Option, error)
}
