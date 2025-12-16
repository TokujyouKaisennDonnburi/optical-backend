package query

import "context"

type SuggestOptionsInput struct {
	Message string
}

type SuggestOptionsOutput struct {
	Id   int32
	Name string
}

func (q *AgentQuery) SuggestOptions(ctx context.Context, input SuggestOptionsInput) ([]SuggestOptionsOutput, error) {
	options, err := q.optionRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	options, err = q.optionAgentRepository.SuggestOptions(ctx, input.Message, options)
	if err != nil {
		return nil, err
	}
	outputs := make([]SuggestOptionsOutput, len(options))
	for i, option := range options {
		outputs[i] = SuggestOptionsOutput{
			Id:   option.Id,
			Name: option.Name,
		}
	}
	return outputs, nil
}
