package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CalendarModel struct {
	CalendarId         uuid.UUID      `db:"id"`
	CalendarName       string         `db:"name"`
	CalendarOptionId   sql.NullInt64  `db:"option_id"`
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
	`
	calendarModels := []CalendarModel{}
	fmt.Println(calendarId.String())
	err := tx.SelectContext(ctx, &calendarModels, query, calendarId)
	if err != nil {
		return nil, err
	}
	if len(calendarModels) == 0 {
		return nil, errors.New("calendar not found")
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
			Id:   int(calendarModel.CalendarOptionId.Int64),
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
