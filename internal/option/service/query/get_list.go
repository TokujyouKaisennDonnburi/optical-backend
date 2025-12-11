package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type ListQueryInput struct {
	UserId uuid.UUID
}

// option 一覧取得
func (q *OptionQuery) GetListOption(ctx context.Context, input ListQueryInput) ([]option.Option, error) {
	outputs, err := q.optionRepository.GetList(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	result := make([]option.Option, len(outputs))
	for i, o := range outputs {
		result[i] = option.Option{
			Id:         o.Id,
			Name:       o.Name,
			Deprecated: o.Deprecated,
		}
	}
	return result, nil 
}
