package gateway

import "google.golang.org/genai"

type OptionAgentGeminiRepository struct {
	client *genai.Client
}

type OptionAgentOpenRouterRepository struct {
	apiKey string
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
	apiKey string,
) *OptionAgentOpenRouterRepository {
	if apiKey == "" {
		panic("OpenRouter API KEY is nil")
	}
	return &OptionAgentOpenRouterRepository{
		apiKey: apiKey,
	}
}
