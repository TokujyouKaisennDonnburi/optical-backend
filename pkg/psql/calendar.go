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

// ユーザーIDとカレンダーIDからカレンダーを取得する
func FindCalendarByUserIdAndId(ctx context.Context, tx *sqlx.Tx, userId, calendarId uuid.UUID) (*calendar.Calendar, error) {
	query := `
		SELECT calendars.id, image_id, name, color
		FROM calendars
		JOIN calendar_members ON calendar_members.calendar_id = calendars.id
		WHERE calendars.id = $1
		AND calendar_members.user_id = $2
		AND calendar_members.joined_at IS NOT NULL
		AND calendars.deleted_at IS NULL
	`

	var calendarModel CalendarModel
	err := tx.GetContext(ctx, &calendarModel, query, calendarId, userId)
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
