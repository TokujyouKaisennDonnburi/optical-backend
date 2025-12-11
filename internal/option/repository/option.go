package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type OptionRepository interface {
	FindByIds(ctx context.Context, ids []int) ([]option.Option, error)
	List(ctx context.Context, userId uuid.UUID) ([]option.Option, error)
}
