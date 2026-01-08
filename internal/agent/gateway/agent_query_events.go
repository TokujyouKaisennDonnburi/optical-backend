package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type EventAndCalendarModel struct {
	CalendarId    uuid.UUID `db:"calendar_id"`
	CalendarName  string    `db:"calendar_name"`
	CalendarColor string    `db:"calendar_color"`
	EventId       uuid.UUID `db:"event_id"`
	EventTitle    string    `db:"event_title"`
	Location      string    `db:"location"`
	Memo          string    `db:"memo"`
	IsAllday      bool      `db:"all_day"`
	StartAt       time.Time `db:"start_at"`
	EndAt         time.Time `db:"end_at"`
}

func (r *AgentQueryPsqlRepository) FindEventByUserIdAndDate(
	ctx context.Context,
	userId uuid.UUID,
	startAt, endAt time.Time,
) ([]agent.AnalyzableEvent, error) {
	events := []agent.AnalyzableEvent{}
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		var models []EventAndCalendarModel
		query := `
			SELECT 
				calendars.id AS calendar_id, calendars.name AS calendar_name, calendars.color AS calendar_color,
				events.id AS event_id, events.title AS event_title, location, memo, start_at, end_at, all_day
			FROM events
			JOIN event_locations
				ON events.id = event_locations.event_id
			JOIN calendars
				ON events.calendar_id = calendars.id
			JOIN calendar_members
				ON calendar_members.calendar_id = events.calendar_id
			WHERE
				events.deleted_at IS NULL 
				AND calendars.deleted_at IS NULL 
				AND	calendar_members.user_id = $1
				AND	
				(
					events.start_at >= $2 AND events.start_at <= $3
					OR events.end_at <= $3 AND events.end_at >= $2
					OR events.start_at <= $2 AND events.end_at >= $3
				)
		`
		err := tx.SelectContext(ctx, &models, query, userId, startAt, endAt)
		if err != nil {
			logrus.WithError(err).Error("events query error")
			return err
		}
		for _, output := range models {
			events = append(events, agent.AnalyzableEvent{
				CalendarId:    output.CalendarId.String(),
				CalendarName:  output.CalendarName,
				CalendarColor: output.CalendarColor,
				Id:            output.EventId.String(),
				Title:         output.EventTitle,
				Location:      output.Location,
				Memo:          output.Memo,
				StartAt:       output.StartAt,
				EndAt:         output.EndAt,
				IsAllday:      output.IsAllday,
			})
		}
		return nil
	})
	return events, err
}

func (r *AgentQueryPsqlRepository) FindCalendarEventByUserIdAndDate(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	startAt, endAt time.Time,
) ([]agent.AnalyzableEvent, error) {
	events := []agent.AnalyzableEvent{}
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		var models []EventAndCalendarModel
		query := `
			SELECT 
				calendars.id AS calendar_id, calendars.name AS calendar_name, calendars.color AS calendar_color,
				events.id AS event_id, events.title AS event_title, location, memo, start_at, end_at, all_day
			FROM events
			JOIN event_locations
				ON events.id = event_locations.event_id
			JOIN calendars
				ON events.calendar_id = calendars.id
			JOIN calendar_members
				ON calendar_members.calendar_id = events.calendar_id
			WHERE
				events.deleted_at IS NULL 
				AND calendars.id = $2
				AND calendars.deleted_at IS NULL 
				AND	calendar_members.user_id = $1
				AND	
				(
					events.start_at >= $3 AND events.start_at <= $4
					OR events.end_at <= $4 AND events.end_at >= $3
					OR events.start_at <= $3 AND events.end_at >= $4
				)
		`
		err := tx.SelectContext(ctx, &models, query, userId, calendarId, startAt, endAt)
		if err != nil {
			logrus.WithError(err).Error("events query error")
			return err
		}
		for _, output := range models {
			events = append(events, agent.AnalyzableEvent{
				CalendarId:    output.CalendarId.String(),
				CalendarName:  output.CalendarName,
				CalendarColor: output.CalendarColor,
				Id:            output.EventId.String(),
				Title:         output.EventTitle,
				Location:      output.Location,
				Memo:          output.Memo,
				StartAt:       output.StartAt,
				EndAt:         output.EndAt,
				IsAllday:      output.IsAllday,
			})
		}
		return nil
	})
	return events, err
}
