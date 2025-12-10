package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)

type OptionRepository interface {
	FindByIds(ctx context.Context, ids []int) ([]option.Option, error)
}
