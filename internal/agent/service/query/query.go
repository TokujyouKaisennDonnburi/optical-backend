package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
)

type AgentQuery struct {
	agentRepository       repository.AgentRepository
	optionRepository      optionRepo.OptionRepository
	eventRepository       calendarRepo.EventRepository
	optionAgentRepository repository.OptionAgentRepository
}

func NewAgentQuery(
	agentRepository repository.AgentRepository,
	optionRepository optionRepo.OptionRepository,
	eventRepository calendarRepo.EventRepository,
	optionAgentRepository repository.OptionAgentRepository,
) *AgentQuery {
	if agentRepository == nil {
		panic("agentRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	if eventRepository == nil {
		panic("eventRepository is nil")
	}
	if optionAgentRepository == nil {
		panic("optionAgentRepository is nil")
	}
	return &AgentQuery{
		agentRepository:       agentRepository,
		optionRepository:      optionRepository,
		eventRepository:       eventRepository,
		optionAgentRepository: optionAgentRepository,
	}
}
