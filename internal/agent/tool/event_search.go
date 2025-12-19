package tool

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type EventSearchTool struct {
	db     *sqlx.DB
	userId uuid.UUID
}

func NewEventSearchTool(db *sqlx.DB, userId uuid.UUID) *EventSearchTool {
	if db == nil {
		panic("db is nil")
	}
	return &EventSearchTool{
		db:     db,
		userId: userId,
	}
}

func (t EventSearchTool) Name() string {
	return "予定検索ツール"
}

func (t EventSearchTool) Description() string {
	return "期間を指定して予定を検索します。開始日時と終了日時をRFC3339形式で指定して検索します"
}

func (t EventSearchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"start_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
			"end_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
		},
		"required": []string{"start_at", "end_at"},
	}
}

type EventSearchInput struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

func (t EventSearchTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Debug("event search called")
	var inputModel EventSearchInput
	if err := json.Unmarshal([]byte(input), &inputModel); err != nil {
		return "", err
	}
	if inputModel.StartAt.IsZero() || inputModel.EndAt.IsZero() {
		logrus.WithField("user_input", input).Error("invalid user input time")
		return "", errors.New("input time is nil")
	}
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
				events.start_at > $2 AND events.end_at < $3
			)
	`
	err := t.db.SelectContext(ctx, &models, query, t.userId, inputModel.StartAt, inputModel.EndAt)
	if err != nil {
		logrus.WithError(err).Error("events query error")
		return "", err
	}
	output, err := json.Marshal(models)
	if err != nil {
		logrus.WithError(err).Error("events unmarshal error")
		return "", err
	}
	logrus.WithFields(logrus.Fields{
		"len":      len(models),
		"start_at": inputModel.StartAt.Format("2006-01-02 15:04:05"),
		"end_at":   inputModel.EndAt.Format("2006-01-02 15:04:05"),
	}).Debug("event search tool called")
	return string(output), nil
}

type EventAndCalendarModel struct {
	CalendarId    uuid.UUID `json:"calendar_id" db:"calendar_id"`
	CalendarName  string    `json:"calendar_name" db:"calendar_name"`
	CalendarColor string    `json:"calendar_color" db:"calendar_color"`
	EventId       uuid.UUID `json:"event_id" db:"event_id"`
	EventTitle    string    `json:"event_title" db:"event_title"`
	Location      string    `json:"location" db:"location"`
	Memo          string    `json:"memo" db:"memo"`
	IsAllday      bool      `json:"is_allday" db:"all_day"`
	StartAt       time.Time `json:"start_at" db:"start_at"`
	EndAt         time.Time `json:"end_at" db:"end_at"`
}
