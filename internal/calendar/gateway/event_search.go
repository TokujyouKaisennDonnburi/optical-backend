package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/google/uuid"
)

// 検索結果用のクエリモデル
// タイトル、メモ、場所を検索
type EventSearchQueryModel struct {
	CalendarId    uuid.UUID `db:"calendar_id"`
	CalendarName  string    `db:"calendar_name"`
	CalendarColor string    `db:"calendar_color"`
	EventId       uuid.UUID `db:"event_id"`
	EventTitle    string    `db:"event_title"`
	Location      string    `db:"location"`
	Memo          string    `db:"memo"`
	StartAt       time.Time `db:"start_at"`
	EndAt         time.Time `db:"end_at"`
	IsAllDay      bool      `db:"all_day"`
}

// 検索ロジック
const eventSearchBaseSQL = `
    FROM events
    LEFT JOIN event_locations ON events.id = event_locations.event_id
    JOIN calendars ON events.calendar_id = calendars.id
    JOIN calendar_members ON calendar_members.calendar_id = events.calendar_id
    WHERE
        events.deleted_at IS NULL
        AND calendars.deleted_at IS NULL
        AND calendar_members.user_id = :user_id
        AND calendar_members.joined_at IS NOT NULL
        AND events.start_at >= :start_from
        AND events.start_at <= :start_to
        AND (
	    events @@@ paradedb.parse(:query)
            OR
            COALESCE(event_locations.location, '') ILIKE '%' || :query || '%'
        )
	`

// pg_searchを使ったイベント検索
func (r *EventPsqlRepository) SearchEvents(
	ctx context.Context,
	params repository.SearchEventsParams,
) (*output.EventSearchQueryOutput, error) {
	// デフォルト値適用
	params.ApplyDefaults()

	// 検索クエリ（pg_search版）
	query := `
    SELECT
        calendars.id AS calendar_id,
        calendars.name AS calendar_name,
        calendars.color AS calendar_color,
        events.id AS event_id,
        events.title AS event_title,
        COALESCE(event_locations.location, '') AS location,
        COALESCE(events.memo, '') AS memo,
        events.start_at,
        events.end_at,
        events.all_day
	` + eventSearchBaseSQL + `
       ORDER BY events.start_at DESC
    LIMIT :limit OFFSET :offset
`

	rows, err := r.db.NamedQueryContext(ctx, query, map[string]any{
		"user_id":    params.UserId,
		"query":      params.Query,
		"start_from": params.StartFrom,
		"start_to":   params.StartTo,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []EventSearchQueryModel
	for rows.Next() {
		var model EventSearchQueryModel
		if err := rows.StructScan(&model); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 総件数取得
	totalCount, err := r.countSearchEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	// 出力形式に変換
	items := make([]output.EventSearchQueryOutputItem, len(models))
	for i, model := range models {
		items[i] = output.EventSearchQueryOutputItem{
			CalendarId:    model.CalendarId,
			CalendarName:  model.CalendarName,
			CalendarColor: model.CalendarColor,
			EventId:       model.EventId,
			EventTitle:    model.EventTitle,
			Location:      model.Location,
			Memo:          model.Memo,
			StartAt:       model.StartAt,
			EndAt:         model.EndAt,
			IsAllDay:      model.IsAllDay,
		}
	}

	return &output.EventSearchQueryOutput{
		Items: items,
		Total: totalCount,
		Limit: params.Limit,
	}, nil
}

// 検索結果の総件数を取得
func (r *EventPsqlRepository) countSearchEvents(
	ctx context.Context,
	params repository.SearchEventsParams,
) (int, error) {
	query := `SELECT COUNT(DISTINCT events.id) ` + eventSearchBaseSQL

	rows, err := r.db.NamedQueryContext(ctx, query, map[string]any{
		"user_id":    params.UserId,
		"query":      params.Query,
		"start_from": params.StartFrom,
		"start_to":   params.StartTo,
	})
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}
