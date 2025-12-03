package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)


type CalendarPsqlRepository struct {
	db *sqlx.DB
}

func NewCalendarPsqlRepository(db *sqlx.DB) *CalendarPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &CalendarPsqlRepository{
		db: db,
	}
}

// スケジュールを新規作成する
func (r *CalendarPsqlRepository) Create(
	ctx context.Context,
	userId, imageId uuid.UUID,
	optionIds []uuid.UUID,
	createFn func(user *user.User, image *calendar.Image, options []option.Option) (*calendar.Calendar, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// オプション取得
		options, err := psql.FindOptionsByIds(ctx, tx, optionIds)
		if err != nil {
			return err
		}
		// ユーザー取得
		user, err := psql.FindUserById(ctx, tx, userId)
		if err != nil {
			return err
		}
		// 画像を取得
		image, err := psql.FindImageById(ctx, tx, imageId)
		if err != nil {
			return err
		}
		// スケジュール作成関数を実行
		calendar, err := createFn(user, image, options)
		if err != nil {
			return err
		}
		// スケジュール作成
		query := `
			INSERT INTO calendars(id, name, color, image_id)
			VALUES (:id, :name, :color, :imageId)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":    calendar.Id,
			"name":  calendar.Name,
			"color": calendar.Color,
			"imageId": uuid.NullUUID{
				UUID:  image.Id,
				Valid: image.Valid,
			},
		})
		if err != nil {
			return err
		}
		// メンバー作成
		query = `
			INSERT INTO calendar_members(calendar_id, user_id, joined_at)
			VALUES (:calendarId, :userId, :joinedAt)
		`
		calendarMemberMaps := []map[string]any{}
		for _, member := range calendar.Members {
			calendarMemberMaps = append(calendarMemberMaps, map[string]any{
				"calendarId": calendar.Id,
				"userId":     member.UserId,
				"joinedAt":   member.JoinedAt,
			})
		}
		_, err = tx.NamedExecContext(ctx, query, calendarMemberMaps)
		if err != nil {
			return err
		}
		// オプション設定
		if len(options) > 0 {
			query = `
				INSERT INTO calendar_options(calendar_id, option_id)
				VALUES (:calendarId, :optionId)
			`
			calendarOptionMaps := []map[string]any{}
			for _, option := range calendar.Options {
				calendarOptionMaps = append(calendarOptionMaps, map[string]any{
					"calendarId": calendar.Id,
					"optionId":   option.Id,
				})
			}
			_, err = tx.NamedExecContext(ctx, query, calendarOptionMaps)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type CalendarListQueryModel struct{
	Id uuid.UUID `db:"id"`
	Name  string `db:"name"`
	Color string `db:"color"`
}

// ユーザーが所属するカレンダー一覧を取得する
func (r *CalendarPsqlRepository) FindByUserId(ctx context.Context, userId uuid.UUID) ([]output.CalendarQueryOutput, error) {
	query := `
		SELECT c.id, c.name, c.color
		FROM calendars c
		INNER JOIN calendar_members m ON c.id = m.calendar_id
		WHERE m.user_id = $1
		AND c.deleted_at IS NULL
		ORDER BY c.id
	`
	var rows []CalendarListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, err
	}

	calendars := make([]output.CalendarQueryOutput, len(rows))
	for i, row := range rows {
		calendars[i] = output.CalendarQueryOutput{
			Id:    row.Id,
			Name:  row.Name,
			Color: row.Color,
		}
	}
	return calendars, nil
}

type CalendarQueryModel struct {
	Id      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	Color   string     `db:"color"`
	Image   calendar.Image    `db:"image"`
	Members []calendar.Member `db:"member"`
	Options []option.Option   `db:"option"`
}

// calendarの単体取得
func (r *CalendarPsqlRepository) FindById(ctx context.Context, id uuid.UUID) (*calendar.Calendar, error) {
	query := `
        SELECT c.id, c.name, c.color, i.id AS image_id, i.url, m.user_id, o.option_id
        FROM calendars c
		LEFT JOIN calendar_images i ON i.id = c.image_id
		INNER JOIN calendar_options o ON o.calendar_id = c.id
		INNER JOIN calendar_members m ON m.calendar_id = c.id
        WHERE c.id = $1
		ORDER BY c.id
    `
	model := CalendarQueryModel{}
	err := r.db.SelectContext(ctx, &model, query, id)
	if err != nil {
		return nil, err
	}
	return &calendar.Calendar{
		Id:      model.Id,
		Name:    model.Name,
		Color:   model.Color,
		Image:   model.Image,
		Members: model.Members,
		Options: model.Options,
	},nil
}

