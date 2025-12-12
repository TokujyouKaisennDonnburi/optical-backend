package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type OptionRepository interface {
	FindByIds(ctx context.Context, ids []int32) ([]option.Option, error)
	GetList(ctx context.Context, userId uuid.UUID) ([]option.Option, error)
}
