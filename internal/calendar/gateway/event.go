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

type EventTodayQueryModel struct {
	CalendarId    uuid.UUID `db:"calendar_id"`
	CalendarName  string    `db:"calendar_name"`
	CalendarColor string    `db:"calendar_color"`
	EventId       uuid.UUID `db:"id"`
	EventTitle    string    `db:"title"`
	EventColor    string    `db:"event_color"`
	Location      string    `db:"location"`
	Memo          string    `db:"memo"`
	StartAt       time.Time `db:"start_at"`
	EndAt         time.Time `db:"end_at"`
	IsAllDay      bool      `db:"all_day"`
}

func (r *EventPsqlRepository) GetEventsByDate(
	ctx context.Context,
	userId uuid.UUID,
	datetime time.Time,
) ([]output.EventTodayQueryOutputItem, error) {
	query := `
		SELECT 
			calendar_id, calendar_name, calendar_color,
			event_id, event_title, event_color, location, memo, start_at, end_at, all_day
		FROM events
		JOIN calendars
			ON events.calendar_id = calendars.id
		JOIN calendar_members
			ON calendar_members.calendar_id = events.calendar_id
		WHERE
			calendar_member.user_id = $1
			AND	
			(
				events.start_at::date = $2 AND events.end_at::date < $2 
				OR events.start_at::date > $2 AND events.end_at::date > $2
			)
	`
	var models []EventTodayQueryModel
	err := r.db.SelectContext(ctx, &models, query, userId, datetime.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	var outputs []output.EventTodayQueryOutputItem
	for _, model := range models {
		outputs = append(outputs, output.EventTodayQueryOutputItem{
			CalendarId:    model.CalendarId,
			CalendarName:  model.CalendarName,
			CalendarColor: model.CalendarColor,
			Id:            model.EventId,
			Title:         model.EventTitle,
			Color:         model.EventColor,
			Location:      model.Location,
			Memo:          model.Memo,
			StartAt:       model.StartAt,
			EndAt:         model.EndAt,
			IsAllDay:      model.IsAllDay,
		})
	}
	return outputs, nil
}
