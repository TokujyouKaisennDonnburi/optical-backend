package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CalendarModel struct {
	CalendarId         uuid.UUID      `db:"id"`
	CalendarName       string         `db:"name"`
	CalendarOptionId   sql.NullInt32  `db:"option_id"`
	CalendarOptionName sql.NullString `db:"option_name"`
}

type MemberModel struct {
	CalendarId uuid.UUID    `db:"calendar_id"`
	UserId     uuid.UUID    `db:"user_id"`
	UserName   string       `db:"name"`
	JoinedAt   sql.NullTime `db:"joined_at"`
}

// トランザクションでIDからカレンダーを取得
func FindCalendarById(ctx context.Context, tx *sqlx.Tx, calendarId uuid.UUID) (*calendar.Calendar, error) {
	query := `
		SELECT 
			calendars.id, calendars.name,
			calendar_options.option_id AS option_id, 
			options.name AS option_name
		FROM calendars
		LEFT JOIN calendar_options
			ON calendars.id = calendar_options.calendar_id
		LEFT JOIN options
			ON options.id = calendar_options.option_id
		WHERE 
			calendars.id = $1 AND calendars.deleted_at IS NULL
		ORDER BY calendars.id
	`
	calendarModels := []CalendarModel{}
	err := tx.SelectContext(ctx, &calendarModels, query, calendarId)
	if err != nil {
		return nil, err
	}
	if len(calendarModels) == 0 {
		return nil, apperr.NotFoundError("calendar not found")
	}
	query = `
		SELECT 
			calendar_members.calendar_id, calendar_members.user_id, calendar_members.joined_at,
			users.name
		FROM calendar_members
		JOIN users
			ON calendar_members.user_id = users.id
		WHERE 
			calendar_members.calendar_id = $1
	`
	memberModels := []MemberModel{}
	err = tx.SelectContext(ctx, &memberModels, query, calendarId)
	if err != nil {
		return nil, err
	}
	return modelsToCalendar(calendarModels, memberModels)
}

type CalendarAndImageModel struct {
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

	var calendarModel CalendarAndImageModel
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

// カレンダーとメンバーのモデルをカレンダーに変換
func modelsToCalendar(calendarModels []CalendarModel, memberModels []MemberModel) (*calendar.Calendar, error) {
	if len(calendarModels) == 0 {
		return nil, errors.New("calendar is empty")
	}
	options := []option.Option{}
	for _, calendarModel := range calendarModels {
		if !calendarModel.CalendarOptionId.Valid || !calendarModel.CalendarOptionName.Valid {
			continue
		}
		options = append(options, option.Option{
			Id:   calendarModel.CalendarOptionId.Int32,
			Name: calendarModel.CalendarOptionName.String,
		})
	}
	members := []calendar.Member{}
	for _, memberModel := range memberModels {
		members = append(members, calendar.Member{
			UserId: memberModel.UserId,
			Name:   memberModel.UserName,
		})
	}
	return &calendar.Calendar{
		Id:      calendarModels[0].CalendarId,
		Name:    calendarModels[0].CalendarName,
		Members: members,
		Options: options,
	}, nil
}

func IsUserInCalendarMembers(ctx context.Context, tx *sqlx.Tx, userId, calendarId uuid.UUID) (bool, error) {
	exists := false
	query := `
		SELECT 1
		FROM calendar_members
		WHERE calendar_members.calendar_id = $2
			AND calendar_members.user_id = $1
	`
	err := tx.GetContext(ctx, &exists, query, userId, calendarId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, apperr.ForbiddenError(err.Error())
		}
		return false, err
	}
	return exists, nil
}
