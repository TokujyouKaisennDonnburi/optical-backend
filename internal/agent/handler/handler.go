package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/query"
)

type AgentHandler struct {
	agentQuery   *query.AgentQuery
	agentCommand *command.AgentCommand
}

func NewAgentHandler(
	agentQuery *query.AgentQuery,
	agentCommand *command.AgentCommand,
) *AgentHandler {
	if agentQuery == nil {
		panic("agentQuery is nil")
	}
	if agentCommand == nil {
		panic("agentCommand is nil")
	}
	return &AgentHandler{
		agentQuery:   agentQuery,
		agentCommand: agentCommand,
	}
}
