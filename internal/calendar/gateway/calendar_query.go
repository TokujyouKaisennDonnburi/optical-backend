package gateway

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
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
	`
	calendarModels := []CalendarModel{}
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

type CalendarListQueryModel struct {
	Id       uuid.UUID      `db:"id"`
	Name     string         `db:"name"`
	Color    string         `db:"color"`
	ImageId  uuid.NullUUID  `db:"image_id"`
	ImageUrl sql.NullString `db:"image_url"`
}

// ユーザーが所属するカレンダー一覧を取得する
func (r *CalendarPsqlRepository) FindByUserId(ctx context.Context, userId uuid.UUID) ([]output.CalendarListQueryOutput, error) {
	query := `
		SELECT 
			c.id, c.name, c.color, c.image_id, ci.url AS image_url
		FROM calendars c
		INNER JOIN calendar_members m 
			ON c.id = m.calendar_id
		LEFT JOIN calendar_images ci
			ON c.image_id = ci.id
		WHERE 
			m.user_id = $1
			AND c.deleted_at IS NULL
		ORDER BY c.id DESC
	`
	var rows []CalendarListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, err
	}

	calendars := make([]output.CalendarListQueryOutput, len(rows))
	for i, row := range rows {
		calendars[i] = output.CalendarListQueryOutput{
			Id:    row.Id,
			Name:  row.Name,
			Color: row.Color,
			Image: calendar.Image{
				Id:    row.ImageId.UUID,
				Url:   row.ImageUrl.String,
				Valid: row.ImageId.Valid && row.ImageUrl.Valid,
			},
		}
	}
	return calendars, nil
}

type CalendarQueryModel struct {
	Id       uuid.UUID         `db:"id"`
	Name     string            `db:"name"`
	Color    calendar.Color    `db:"color"`
	ImageId  uuid.NullUUID     `db:"imageId"`
	ImageUrl sql.NullString    `db:"imageUrl"`
	Members  []calendar.Member `db:"member"`
	Options  []option.Option   `db:"option"`
}
type CalendarRow struct {
	Id       uuid.UUID      `db:"id"`
	Name     string         `db:"name"`
	Color    calendar.Color `db:"color"`
	ImageId  uuid.NullUUID  `db:"image_id"`
	ImageUrl sql.NullString `db:"image_url"`
}

// calendar単体取得
func (r *CalendarPsqlRepository) FindByUserCalendarId(ctx context.Context, userId, calendarId uuid.UUID) (*calendar.Calendar, error) {
	// calendar & image
	query := `
	SELECT
	calendars.id,
	calendars.name,
	calendars.color,
	calendars.image_id,
	calendar_images.url AS image_url
	FROM calendars
	LEFT JOIN calendar_images ON
	calendar_images.id = calendars.image_id
	WHERE calendars.id = $1
	`
	var calRow CalendarRow
	err := r.db.GetContext(ctx, &calRow, query, calendarId)
	if err != nil {
		return nil, err
	}
	query = `
	SELECT
	calendar_members.user_id,
	users.name,
	calendar_members.joined_at
	FROM calendar_members
	INNER JOIN users ON
	users.id = calendar_members.user_id
	WHERE calendar_members.calendar_id = $1
	`
	var memberRows []calendar.Member
	err = r.db.SelectContext(ctx, &memberRows, query, calendarId)
	if err != nil {
		return nil, err
	}

	members := make([]calendar.Member, len(memberRows))
	for i, row := range memberRows {
		members[i] = calendar.Member{
			UserId:   row.UserId,
			Name:     row.Name,
			JoinedAt: row.JoinedAt,
		}
	}
	query = `
	SELECT
	calendar_options.option_id AS id,
	options.name,
	options.deprecated
	FROM calendar_options
	INNER JOIN options ON
	options.id = calendar_options.option_id
	WHERE calendar_options.calendar_id = $1
	`
	var optionRows []option.Option
	err = r.db.SelectContext(ctx, &optionRows, query, calendarId)
	if err != nil {
		return nil, err
	}
	options := make([]option.Option, len(optionRows))
	for i, row := range optionRows {
		options[i] = option.Option{
			Id:         row.Id,
			Name:       row.Name,
			Deprecated: row.Deprecated,
		}
	}
	// bind
	return &calendar.Calendar{
		Id:    calRow.Id,
		Name:  calRow.Name,
		Color: calRow.Color,
		Image: calendar.Image{
			Id:    calRow.ImageId.UUID,
			Url:   calRow.ImageUrl.String,
			Valid: calRow.ImageId.Valid,
		},
		Members: members,
		Options: options,
	}, nil
}
