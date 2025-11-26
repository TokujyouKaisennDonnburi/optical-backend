package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
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
	userId uuid.UUID,
	optionIds []uuid.UUID,
	createFn func(user *user.User, options []option.Option) (*calendar.Calendar, error),
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
		// スケジュール作成関数を実行
		calendar, err := createFn(user, options)
		if err != nil {
			return err
		}
		// スケジュール作成
		query := `
			INSERT INTO calendars(id, name)
			VALUES (:id, :name)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":   calendar.Id,
			"name": calendar.Name,
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
				calendarMemberMaps = append(calendarMemberMaps, map[string]any{
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
