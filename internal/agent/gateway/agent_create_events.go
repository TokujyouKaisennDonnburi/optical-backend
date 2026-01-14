package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (*AgentCommandPsqlRepository) CreateEvents(
	ctx context.Context,
	userId uuid.UUID,
	calendarId uuid.UUID,
	createFn func(*calendar.Calendar) ([]calendar.Event, error),
) error {
	return transact.Transact(ctx, func(tx *sqlx.Tx) error {
		calendar, err := psql.FindCalendarById(ctx, tx, calendarId)
		if err != nil {
			return err
		}
		events, err := createFn(calendar)
		if err != nil {
			return err
		}
		eventMapList := make([]map[string]any, len(events))
		for i, event := range events {
			eventMapList[i] = map[string]any{
				"id":         event.Id,
				"calendarId": event.CalendarId,
				"title":      event.Title,
				"memo":       event.Memo,
				"color":      event.Color,
				"location":   event.Location,
				"allDay":     event.ScheduledTime.AllDay,
				"startAt":    event.ScheduledTime.StartTime,
				"endAt":      event.ScheduledTime.EndTime,
				"createdAt":  time.Now().UTC(),
				"updatedAt":  time.Now().UTC(),
			}
		}
		query := `
			INSERT INTO events(id, calendar_id, title, memo, color, all_day, start_at, end_at, created_at, updated_at)
			VALUES(:id, :calendarId, :title, :memo, :color, :allDay, :startAt, :endAt, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, eventMapList)
		if err != nil {
			return err
		}
		query = `
			INSERT INTO event_locations(event_id, location)
			VALUES(:id, :location)
		`
		_, err = tx.NamedExecContext(ctx, query, eventMapList)
		return err
	})
}
