package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)

type OptionAgentRepository interface {
	SuggestOptions(ctx context.Context, requset string, options []option.Option) ([]option.Option, error)
}
