package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"


type OptionQuery struct {
	optionRepository repository.OptionRepository
}

func NewOptionQuery(optionRepo repository.OptionRepository) *OptionQuery {
	return &OptionQuery{
		optionRepository: optionRepo,
	}
}
