package gateway

import (
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
)

type AgentOpenRouterRepository struct {
	openRouter *openrouter.OpenRouter
}

type OptionAgentOpenRouterRepository struct {
	openRouter *openrouter.OpenRouter
}

func NewAgentOpenRouterRepository(openRouter *openrouter.OpenRouter) *AgentOpenRouterRepository {
	if openRouter == nil {
		panic("openRouter is nil")
	}
	return &AgentOpenRouterRepository{
		openRouter: openRouter,
	}
}

func NewOptionAgentOpenRouterRepository(
	openRouter *openrouter.OpenRouter,
) *OptionAgentOpenRouterRepository {
	if openRouter == nil {
		panic("openRouter is nil")
	}
	return &OptionAgentOpenRouterRepository{
		openRouter: openRouter,
	}
}
