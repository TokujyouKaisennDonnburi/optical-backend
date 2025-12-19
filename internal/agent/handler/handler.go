package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/query"
)

type AgentHandler struct {
	agentQuery *query.AgentQuery
}

func NewAgentHandler(
	agentQuery *query.AgentQuery,
) *AgentHandler {
	return &AgentHandler{
		agentQuery: agentQuery,
	}
}
