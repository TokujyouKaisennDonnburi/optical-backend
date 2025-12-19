package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
)

type AgentRepository interface {
	AnalyzeSchedules(
		ctx context.Context,
		input string,
		schedules []agent.AnalyzableEvent,
		streamingFn func(ctx context.Context, chunk []byte) error,
	) error
}
