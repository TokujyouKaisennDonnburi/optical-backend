package openrouter

import (
	"errors"
	"os"
)

var (
	ErrToolNotFound = errors.New("tool not found")
)

type OpenRouter struct {
	apiKey         string
	model          string
	temperature    float64
	maxTokens      float64
	responseFormat string
	toolChoice     string
	tools          []Tool
}

type ResponseFormat string

const (
	FORMAT_TEXT        ResponseFormat = "text"
	FORMAT_JSON_OBJECT ResponseFormat = "json_object"
	FORMAT_JSON_SCHEMA ResponseFormat = "json_schema"
)

func NewOpenRouter(apiKey string) *OpenRouter {
	return &OpenRouter{
		apiKey:      apiKey,
		temperature: 0.0,
		model:       os.Getenv("AGENT_MODEL"),
	}
}

func (r *OpenRouter) SetModel(model string) {
	r.model = model
}

func (r OpenRouter) WithTools(tools []Tool) OpenRouter {
	r.tools = tools
	r.toolChoice = "auto"
	return r
}

func (r *OpenRouter) getToolByName(name string) (Tool, error) {
	for _, tool := range r.tools {
		if tool.Name() == name {
			return tool, nil
		}
	}
	return nil, ErrToolNotFound
}

func (r OpenRouter) WithModel(model string) OpenRouter {
	r.model = model
	return r
}

func (r OpenRouter) WithTemperature(temperature float64) OpenRouter {
	r.temperature = temperature
	return r
}

func (r OpenRouter) WithResponseFormat(format string) OpenRouter {
	r.responseFormat = format
	return r
}

func (r ResponseFormat) String() string {
	return string(r)
}

func UserMessage(content string) Message {
	return Message{
		Role:    "user",
		Content: content,
	}
}

func SystemMessage(content string) Message {
	return Message{
		Role:    "system",
		Content: content,
	}
}
