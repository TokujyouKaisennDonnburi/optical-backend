package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type ListQueryInput struct {
	UserId     uuid.UUID
}

// option 一覧取得
func (q *OptionQuery) ListGetEvents(ctx context.Context, input ListQueryInput) ([]option.Option, error) {
	outputs, err := q.optionRepository.List(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	return option.Option{
		Id:         outputs.id,
		Name:       outputs.name,
		Depricated: outputs.Depricated,
	}
}
