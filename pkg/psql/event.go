package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventModel struct {
	Id         uuid.UUID `db:"id"`
	CalendarId uuid.UUID `db:"calendar_id"`
	Title      string    `db:"title"`
	Memo       string    `db:"memo"`
	Location   string    `db:"location"`
	IsAllDay   bool      `db:"all_day"`
	StartAt    time.Time `db:"start_at"`
	EndAt      time.Time `db:"end_at"`
}

func FindEventByUserIdAndId(ctx context.Context, tx *sqlx.Tx, userId, eventId uuid.UUID) (*calendar.Event, error) {
	query := `
		SELECT
			events.id, events.calendar_id, title, memo, event_locations.location, all_day, start_at, end_at
		FROM events
		JOIN calendar_members
			ON calendar_members.calendar_id = events.calendar_id
		JOIN event_locations
			ON event_locations.event_id = events.id
		WHERE
			id = $2 AND
			calendar_members.user_id = $1 AND
			calendar_members.joined_at IS NOT NULL AND
			deleted_at IS NULL
		ORDER BY events.id
	`
	var eventModel EventModel
	err := tx.GetContext(ctx, &eventModel, query, userId, eventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFoundError("event not found")
		}
		return nil, err
	}
	return &calendar.Event{
		Id:         eventModel.Id,
		CalendarId: eventModel.CalendarId,
		Title:      eventModel.Title,
		Memo:       eventModel.Memo,
		Location:   eventModel.Location,
		ScheduledTime: calendar.ScheduledTime{
			AllDay:    eventModel.IsAllDay,
			StartTime: eventModel.StartAt,
			EndTime:   eventModel.EndAt,
		},
	}, nil
}
