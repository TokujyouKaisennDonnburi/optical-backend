package gateway

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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
			AND m.joined_at IS NOT NULL
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
	Id      uuid.UUID      `db:"id"`
	Name    string         `db:"name"`
	Color   calendar.Color `db:"color"`
	ImageId uuid.NullUUID  `db:"imageId"`
}
type CalendarImageMember struct {
	Id                   uuid.UUID      `db:"id"`
	Name                 string         `db:"name"`
	Color                calendar.Color `db:"color"`
	ImageId              uuid.NullUUID  `db:"image_id"`
	ImageUrl             sql.NullString `db:"image_url"`
	UserId               uuid.UUID      `db:"user_id"`
	UserName             string         `db:"user_name"`
	JoinedAt             time.Time      `db:"joined_at"`
	AvatarUrl            sql.NullString `db:"avatar_url"`
	AvatarIsRelativePath sql.NullBool   `db:"avatar_is_relative_path"`
}

type OptionModel struct {
	Id         int32  `db:"id"`
	Name       string `db:"name"`
	Deprecated bool   `db:"deprecated"`
}

// calendar単体取得
func (r *CalendarPsqlRepository) FindByUserCalendarId(ctx context.Context, userId, calendarId uuid.UUID) (*calendar.Calendar, error) {
	// calendar & image & member & users
	query := `
	SELECT
	calendars.id, calendars.name, calendars.color,
	calendars.image_id, calendar_images.url AS image_url,
	calendar_members.user_id, calendar_members.joined_at,
	users.name AS user_name,
	avatars.url AS avatar_url, avatars.is_relative_path AS avatar_is_relative_path
	FROM calendars
	LEFT JOIN calendar_images ON calendar_images.id = calendars.image_id
	INNER JOIN calendar_members ON calendar_members.calendar_id = calendars.id
	INNER JOIN users ON users.id = calendar_members.user_id
	LEFT JOIN user_profiles ON user_profiles.user_id = users.id
	LEFT JOIN avatars ON avatars.id = user_profiles.avatar_id
	WHERE calendars.id = $1
	AND calendar_members.joined_at IS NOT NULL
	AND calendars.deleted_at IS NULL`
	calRow := []CalendarImageMember{}
	err := r.db.SelectContext(ctx, &calRow, query, calendarId)
	if err != nil {
		return nil, err
	}
	if len(calRow) == 0 {
		return nil, errors.New("calendar member is not found")
	}
	exists := false
	members := make([]calendar.Member, len(calRow))
	for i, row := range calRow {
		if row.UserId == userId {
			exists = true
		}
		avatarUrl := row.AvatarUrl.String
		if row.AvatarUrl.Valid && row.AvatarIsRelativePath.Bool {
			avatarUrl = storage.GetImageStorageBaseUrl() + "/" + avatarUrl
		}
		members[i] = calendar.Member{
			UserId:    row.UserId,
			Name:      row.UserName,
			JoinedAt:  row.JoinedAt,
			AvatarUrl: avatarUrl,
		}
	}
	if !exists {
		return nil, apperr.ForbiddenError("user is not a member of this calendar")
	}
	// option
	query = `
	SELECT id, name, deprecated 
	FROM options
	WHERE options.id IN (
		SELECT option_id
		FROM calendar_options
		WHERE calendar_id = $1
	);
	`
	optionModels := []OptionModel{}
	err = r.db.SelectContext(ctx, &optionModels, query, calendarId)
	if err != nil {
		return nil, err
	}
	options := make([]option.Option, len(optionModels))
	for i, row := range optionModels {
		options[i] = option.Option{
			Id:         row.Id,
			Name:       row.Name,
			Deprecated: row.Deprecated,
		}
		members := make([]calendar.Member, len(calRow))
		for i, row := range calRow {
			members[i] = calendar.Member{
				UserId:   row.UserId,
				Name:     row.UserName,
				JoinedAt: row.JoinedAt,
			}
		}
		// option
		query = `
		SELECT id, name, deprecated 
		FROM options
		WHERE options.id IN (
			SELECT option_id
			FROM calendar_options
			WHERE calendar_id = $1
		);
		`
		optionModels := []OptionModel{}
		err = r.db.SelectContext(ctx, &optionModels, query, calendarId)
		if err != nil {
			return err
		}
		options := make([]option.Option, len(optionModels))
		for i, row := range optionModels {
			options[i] = option.Option{
				Id:         row.Id,
				Name:       row.Name,
				Deprecated: row.Deprecated,
			}
		}
		// bind
		cal = &calendar.Calendar{
			Id:    calRow[0].Id,
			Name:  calRow[0].Name,
			Color: calRow[0].Color,
			Image: calendar.Image{
				Id:    calRow[0].ImageId.UUID,
				Url:   calRow[0].ImageUrl.String,
				Valid: calRow[0].ImageId.Valid,
			},
			Members: members,
			Options: options,
		}
		return nil
	})
	return cal, err
}

// calendar単体取得
func (r *CalendarPsqlRepository) FindById(ctx context.Context, calendarId uuid.UUID) (*calendar.Calendar, error) {
	var cal *calendar.Calendar
	err := db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		// calendar & image & member & users
		query := `
		SELECT
			calendars.id, calendars.name, calendars.color,
			calendars.image_id, calendar_images.url AS image_url,
			calendar_members.user_id, calendar_members.joined_at,
			users.name AS user_name
		FROM calendars
		LEFT JOIN calendar_images 
			ON calendar_images.id = calendars.image_id
		LEFT JOIN calendar_members 
			ON calendar_members.calendar_id = calendars.id
		INNER JOIN users 
			ON users.id = calendar_members.user_id
		WHERE 
			calendars.id = $1
			AND calendar_members.joined_at IS NOT NULL
			AND calendars.deleted_at IS NULL `
		calRow := []CalendarImageMember{}
		err := r.db.SelectContext(ctx, &calRow, query, calendarId)
		if err != nil {
			return err
		}
		if len(calRow) == 0 {
			return errors.New("calendar member is not found")
		}
		members := make([]calendar.Member, len(calRow))
		for i, row := range calRow {
			members[i] = calendar.Member{
				UserId:   row.UserId,
				Name:     row.UserName,
				JoinedAt: row.JoinedAt,
			}
		}
		// option
		query = `
		SELECT id, name, deprecated 
		FROM options
		WHERE options.id IN (
			SELECT option_id
			FROM calendar_options
			WHERE calendar_id = $1
		);
		`
		optionModels := []OptionModel{}
		err = r.db.SelectContext(ctx, &optionModels, query, calendarId)
		if err != nil {
			return err
		}
		options := make([]option.Option, len(optionModels))
		for i, row := range optionModels {
			options[i] = option.Option{
				Id:         row.Id,
				Name:       row.Name,
				Deprecated: row.Deprecated,
			}
		}
		// bind
		cal = &calendar.Calendar{
			Id:    calRow[0].Id,
			Name:  calRow[0].Name,
			Color: calRow[0].Color,
			Image: calendar.Image{
				Id:    calRow[0].ImageId.UUID,
				Url:   calRow[0].ImageUrl.String,
				Valid: calRow[0].ImageId.Valid,
			},
			Members: members,
			Options: options,
		}
		return nil
	})
	return cal, err
}
