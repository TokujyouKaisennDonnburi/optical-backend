package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/transact"
)

type AgentCommand struct {
	openRouter             *openrouter.OpenRouter
	transactor             transact.TransactionProvider
	agentQueryRepository   repository.AgentQueryRepository
	agentCommandRepository repository.AgentCommandRepository
}

func NewAgentCommand(
	openRouter *openrouter.OpenRouter,
	transactor transact.TransactionProvider,
	agentQueryRepository repository.AgentQueryRepository,
	agentCommandRepository repository.AgentCommandRepository,
) *AgentCommand {
	if openRouter == nil {
		panic("openRouter is nil")
	}
	if transactor == nil {
		panic("transactor is nil")
	}
	if agentQueryRepository == nil {
		panic("agentQueryRepository is nil")
	}
	if agentCommandRepository == nil {
		panic("agentCommandRepository is nil")
	}
	return &AgentCommand{
		openRouter:             openRouter,
		transactor:             transactor,
		agentQueryRepository:   agentQueryRepository,
		agentCommandRepository: agentCommandRepository,
	}
}

