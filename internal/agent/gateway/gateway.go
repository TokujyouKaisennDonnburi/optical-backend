package gateway

import (
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/jmoiron/sqlx"
)

type AgentQueryPsqlRepository struct {
	db *sqlx.DB
}

type AgentCommandPsqlRepository struct {
	db *sqlx.DB
}

type OptionAgentOpenRouterRepository struct {
	openRouter *openrouter.OpenRouter
}

func NewAgentQueryPsqlRepository(
	db *sqlx.DB,
) *AgentQueryPsqlRepository {
	return &AgentQueryPsqlRepository{
		db: db,
	}
}

func NewAgentCommandPsqlRepository(
	db *sqlx.DB,
) *AgentCommandPsqlRepository {
	return &AgentCommandPsqlRepository{
		db: db,
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
