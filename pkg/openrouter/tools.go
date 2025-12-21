package openrouter

import (
	"context"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]any
	Strict() bool
	Call(context.Context, string) (string, error)
}

type FunctionCall struct {
	Id       string `json:"id"`
	Index    int    `json:"index"`
	Type     string `json:"string"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func toolsToRequests(toolList []Tool) []ToolRequest {
	if len(toolList) == 0 {
		return []ToolRequest{}
	}
	tools := make([]ToolRequest, len(toolList))
	for i, t := range toolList {
		tools[i] = ToolRequest{
			Type: "function",
			Function: struct {
				Name        string         `json:"name"`
				Description string         `json:"description"`
				Parameters  map[string]any `json:"parameters"`
				Strict      bool           `json:"strict,omitempty"`
			}{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Parameters(),
				Strict:      t.Strict(),
			},
		}
	}
	return tools
}
