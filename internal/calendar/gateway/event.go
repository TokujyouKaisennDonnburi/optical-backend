package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventPsqlRepository struct {
	db *sqlx.DB
}

func NewEventPsqlRepository(db *sqlx.DB) *EventPsqlRepository {
	return &EventPsqlRepository{
		db: db,
	}
}

func (r *EventPsqlRepository) Create(
	ctx context.Context,
	calendarId uuid.UUID,
	createFn func(calendar *calendar.Calendar) (*calendar.Event, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		calendar, err := FindCalendarById(ctx, tx, calendarId)
		if err != nil {
			return err
		}
		event, err := createFn(calendar)
		if err != nil {
			return err
		}
		query := `
			INSERT INTO events(id, calendar_id, title, memo, color, all_day, start_at, end_at, created_at, updated_at)
			VALUES(:id, :calendarId, :title, :memo, :color, :allDay, :startAt, :endAt, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         event.Id,
			"calendarId": event.CalendarId,
			"title":      event.Title,
			"memo":       event.Memo,
			"color":      event.Color,
			"allDay":     event.ScheduledTime.AllDay,
			"startAt":    event.ScheduledTime.StartTime,
			"endAt":      event.ScheduledTime.EndTime,
			"createdAt":  time.Now(),
			"updatedAt":  time.Now(),
		})
		return err
	})
}

// 特定のカレンダーに対応するイベント一覧取得
func (r *EventPsqlRepository) ListEventsByCalendarId(
	ctx context.Context,
	calendarId uuid.UUID,
) ([]output.EventQueryOutput, error) {
	query := `
		SELECT e.id, e.calendar_id, e.title, e.memo, e.color, e.all_day, e.start_at, e.end_at, e.created_at
		FROM events e
		WHERE e.calendar_id = $1
			AND e.deleted_at IS NULL -- 論理削除されていないもののみ取得
		ORDER BY e.start_at ASC
	`
	events := []output.EventQueryOutput{}
	err := r.db.SelectContext(ctx, &events, query, calendarId)
	if err != nil {
		return nil, err
	}
	return events, nil
}
