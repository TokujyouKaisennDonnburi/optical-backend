package query

import (
	"context"

	"github.com/google/uuid"
)

type OptionListQueryInput struct {
	UserId uuid.UUID
}

type OptionListQueryOutput struct {
	Id         int32
	Name       string
	Deprecated bool
}

// option 一覧取得
func (q *OptionQuery) GetListOption(ctx context.Context, input OptionListQueryInput) ([]OptionListQueryOutput, error) {
	outputs, err := q.optionRepository.GetList(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	result := make([]OptionListQueryOutput, len(outputs))
	for i, o := range outputs {
		result[i] = OptionListQueryOutput{
			Id:         o.Id,
			Name:       o.Name,
			Deprecated: o.Deprecated,
		}
	}
	return result, nil
}
