package repository

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

type EventRepository interface {
	Create(
		ctx context.Context,
		calendarId uuid.UUID,
		createFn func(*calendar.Calendar) (*calendar.Event, error),
	) error

	Update(
		ctx context.Context,
		userId, eventId uuid.UUID,
		updateFn func(*calendar.Event) (*calendar.Event, error),
	) error

	// 一覧取得
	ListEventsByCalendarId(
		ctx context.Context,
		calendarId uuid.UUID,
	) ([]output.EventQueryOutput, error)

	// カレンダーが指定されたユーザーに属しているかチェック
	ExistsCalendarByUserIdAndCalendarId(
		ctx context.Context,
		userId uuid.UUID,
		calendarId uuid.UUID,
	) (bool, error)
	GetEventsByDate(ctx context.Context, userId uuid.UUID, date time.Time) ([]output.EventTodayQueryOutputItem, error)
	GetEventsByMonth(ctx context.Context, userId uuid.UUID, date time.Time) ([]output.EventTodayQueryOutputItem, error)
	FindAnalyzableEventsByUserId(ctx context.Context, userId uuid.UUID) ([]agent.AnalyzableEvent, error)

	// イベント検索
	SearchEvents(
		ctx context.Context,
		params SearchEventsParams,
	) (*output.EventSearchQueryOutput, error)
}

// イベント検索用パラメータ
type SearchEventsParams struct {
	UserId    uuid.UUID
	Query     string    // クエリパラメータ
	StartFrom time.Time // ゼロ値 = デフォルト（3年前）
	StartTo   time.Time // ゼロ値 = デフォルト（3年後）
	Limit     int
	Offset    int
}

// デフォルト値の適用
func (p *SearchEventsParams) ApplyDefaults() {
	now := time.Now()

	// 検索範囲のデフォルト
	if p.StartFrom.IsZero() {
		p.StartFrom = now.AddDate(-3, 0, 0) // 3年前
	}
	if p.StartTo.IsZero() {
		p.StartTo = now.AddDate(3, 0, 0) // 3年後
	}

	// ページネーションのデフォルト
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
