package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
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
