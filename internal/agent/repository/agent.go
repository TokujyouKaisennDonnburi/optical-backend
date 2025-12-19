package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/google/uuid"
)

type AgentQueryRepository interface {
	FindEventByUserIdAndDate(
		ctx context.Context,
		userId uuid.UUID,
		startAt, endAt time.Time,
	) ([]agent.AnalyzableEvent, error)
	FindCalendarByUserId(
		ctx context.Context,
		userId uuid.UUID,
	) ([]agent.AnalyzableCalendar, error)
}
