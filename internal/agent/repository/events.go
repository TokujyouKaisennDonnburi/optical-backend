package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/google/uuid"
)

type AgentEventRepository interface {
	FindByUserIdAndDate(
		ctx context.Context,
		userId uuid.UUID,
		startAt, endAt time.Time,
	) ([]agent.AnalyzableEvent, error)
}
