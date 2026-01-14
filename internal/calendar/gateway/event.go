package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
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
		calendar, err := psql.FindCalendarById(ctx, tx, calendarId)
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
			"createdAt":  time.Now().UTC(),
			"updatedAt":  time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		query = `
			INSERT INTO event_locations(event_id, location)
			VALUES(:eventId, :location)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"eventId":  event.Id,
			"location": event.Location,
		})
		return err
	})
}

func (r *EventPsqlRepository) Update(
	ctx context.Context,
	userId, eventId uuid.UUID,
	updateFn func(event *calendar.Event) (*calendar.Event, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		event, err := psql.FindEventByUserIdAndId(ctx, tx, userId, eventId)
		if err != nil {
			return err
		}
		event, err = updateFn(event)
		if err != nil {
			return err
		}
		query := `
			UPDATE events SET
				title = :title,
				memo = :memo,
				color = :color,
				all_day = :allDay,
				start_at = :startAt,
				end_at = :endAt,
				updated_at = :updatedAt
			WHERE 
				id = :id AND deleted_at IS NULL
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":        event.Id,
			"title":     event.Title,
			"memo":      event.Memo,
			"color":     event.Color,
			"allDay":    event.ScheduledTime.AllDay,
			"startAt":   event.ScheduledTime.StartTime,
			"endAt":     event.ScheduledTime.EndTime,
			"updatedAt": time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		query = `
		UPDATE event_locations
		SET location = :location
		WHERE event_id = :eventId
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"eventId":  event.Id,
			"location": event.Location,
		})
		return err
	})
}
