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
	Id       uuid.UUID      `db:"id"`
	Name     string         `db:"name"`
	Color    string         `db:"color"`
	ImageID  uuid.NullUUID  `db:"image_id"`
	ImageUrl sql.NullString `db:"image_url"`
}

// ユーザーIDとカレンダーIDからカレンダーを取得する
func FindCalendarByUserIdAndId(ctx context.Context, tx *sqlx.Tx, userId, calendarId uuid.UUID) (*calendar.Calendar, error) {
	query := `
		SELECT calendars.id, image_id, name, color, calendar_images.url AS image_url
		FROM calendars
		LEFT JOIN calendar_images ON calendar_images.id = calendars.image_id
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
		Image: calendar.Image{
			Id:    calendarModel.ImageID.UUID,
			Url:   calendarModel.ImageUrl.String,
			Valid: calendarModel.ImageID.Valid && calendarModel.ImageUrl.Valid,
		},
	}, nil
}
