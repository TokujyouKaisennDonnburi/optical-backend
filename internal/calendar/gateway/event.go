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
		if err != nil {
			return err
		}

		// locationが空でない場合はevent_locationsテーブルにも保存
		if event.Location != "" {
			locationQuery := `
				INSERT INTO event_locations(event_id, location)
				VALUES(:eventId, :location)
			`
			_, err = tx.NamedExecContext(ctx, locationQuery, map[string]any{
				"eventId":  event.Id,
				"location": event.Location,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// 専用モデル
// gateway層専用のモデル
type EventListQueryModel struct {
	Id         uuid.UUID `db:"id"`
	CalendarId uuid.UUID `db:"calendar_id"`
	Title      string    `db:"title"`
	Memo       string    `db:"memo"`
	Color      string    `db:"color"`
	Location   string    `db:"location"`
	AllDay     bool      `db:"all_day"`
	StartAt    string    `db:"start_at"`
	EndAt      string    `db:"end_at"`
	CreatedAt  string    `db:"created_at"`
}

// 特定のカレンダーに対応するイベント一覧取得
func (r *EventPsqlRepository) ListEventsByCalendarId(
	ctx context.Context,
	calendarId uuid.UUID,
) ([]output.EventQueryOutput, error) {
	query := `
		SELECT e.id, e.calendar_id, e.title, e.memo, e.color, COALESCE(el.location, '') as location, e.all_day, e.start_at, e.end_at, e.created_at
		FROM events e
		LEFT JOIN event_locations el ON e.id = el.event_id
		WHERE e.calendar_id = $1
			AND e.deleted_at IS NULL -- 論理削除されていないもののみ取得
		ORDER BY e.start_at ASC
	`

	// クエリ実行
	var rows []EventListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, calendarId)
	if err != nil {
		return nil, err
	}

	// 出力形式に変換
	events := make([]output.EventQueryOutput, len(rows))
	for i, row := range rows {
		events[i] = output.EventQueryOutput{
			Id:         row.Id,
			CalendarId: row.CalendarId,
			Title:      row.Title,
			Memo:       row.Memo,
			Color:      row.Color,
			Location:   row.Location,
			IsAllDay:   row.AllDay,
			StartAt:    row.StartAt,
			EndAt:      row.EndAt,
			CreatedAt:  row.CreatedAt,
		}
	}

	return events, nil
}

// カレンダーが指定されたユーザーに属しているかチェック
func (r *EventPsqlRepository) ExistsCalendarByUserIdAndCalendarId(
	ctx context.Context,
	userId uuid.UUID,
	calendarId uuid.UUID,
) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM calendar_members cm
			JOIN calendars c ON cm.calendar_id = c.id
			WHERE cm.user_id = $1
				AND cm.calendar_id = $2
				AND c.deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userId, calendarId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

type EventTodayQueryModel struct {
	CalendarId    uuid.UUID `db:"calendar_id"`
	CalendarName  string    `db:"calendar_name"`
	CalendarColor string    `db:"calendar_color"`
	EventId       uuid.UUID `db:"event_id"`
	EventTitle    string    `db:"event_title"`
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
			calendars.id AS calendar_id, calendars.name AS calendar_name, calendars.color AS calendar_color,
			events.id AS event_id, events.title AS event_title, events.color AS event_color, location, memo, start_at, end_at, all_day
		FROM events
		JOIN event_locations
			ON events.id = event_locations.event_id
		JOIN calendars
			ON events.calendar_id = calendars.id
		JOIN calendar_members
			ON calendar_members.calendar_id = events.calendar_id
		WHERE
			calendar_members.user_id = $1
			AND	
			(
				events.start_at::date = $2 OR events.end_at::date = $2 
				OR events.start_at::date < $2 AND events.end_at::date > $2
			)
	`
	var models []EventTodayQueryModel
	date := datetime.Format("2006-01-02")
	err := r.db.SelectContext(ctx, &models, query, userId, date)
	if err != nil {
		return nil, err
	}
	outputs := make([]output.EventTodayQueryOutputItem, len(models))
	for i, model := range models {
		outputs[i] = output.EventTodayQueryOutputItem{
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
		}
	}
	return outputs, nil
}
