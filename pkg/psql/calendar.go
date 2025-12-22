package psql

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CalendarModel struct {
	Id    uuid.UUID      `db:"id"`
	Image calendar.Image `db:"image_id"`
	Name  string         `db:"name"`
	Color string         `db:"color"`
}

// カレンダーをカレンダーIDから取得する
func FindCalendarByCalendarId(ctx context.Context, tx *sqlx.Tx, calendarId uuid.UUID) (*calendar.Calendar, error) {
	query := `
		SELECT id, image_id, name, color
		FROM calendars
		WHERE id = $1 AND
		deleted_at IS NULL
		ORDER BY id
	`

	var calendarModel CalendarModel
	err := tx.GetContext(ctx, &calendarModel, query, calendarId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFoundError("calendar not found")
		}
		return nil, err
	}
	return &calendar.Calendar{
		Id:    calendarModel.Id,
		Name:  calendarModel.Name,
		Color: calendar.Color(calendarModel.Color),
		Image: calendarModel.Image,
	}, nil
}
