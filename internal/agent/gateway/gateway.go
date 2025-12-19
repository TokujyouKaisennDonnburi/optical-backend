package gateway

import (
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"google.golang.org/genai"
)

type AgentOpenRouterRepository struct {
	openRouter *openrouter.OpenRouter
}

type OptionAgentGeminiRepository struct {
	client *genai.Client
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

func NewOptionAgentGeminiRepository(
	client *genai.Client,
) *OptionAgentGeminiRepository {
	if client == nil {
		panic("genaiClient is nil")
	}
	return &OptionAgentGeminiRepository{
		client: client,
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
