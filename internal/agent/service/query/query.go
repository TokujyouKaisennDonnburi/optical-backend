package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
)

type AgentQuery struct {
	optionRepository       optionRepo.OptionRepository
	optionAgentRepository repository.OptionAgentRepository
}

func NewAgentQuery(
	optionRepository optionRepo.OptionRepository,
	optionAgentRepository repository.OptionAgentRepository,
) *AgentQuery {
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	if optionAgentRepository == nil {
		panic("optionAgentRepository is nil")
	}
	return &AgentQuery{
		optionRepository:       optionRepository,
		optionAgentRepository: optionAgentRepository,
	}
}
