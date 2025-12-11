package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

type ListQueryInput struct {
	UserId uuid.UUID
	CalendarId uuid.UUID
}

// option 一覧取得
func (q *OptionQuery) ListGetEvents(ctx context.Context, ListQueryInput)([]option.Option, error) {
}

