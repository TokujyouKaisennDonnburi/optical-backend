package gateway

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
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
	imageId uuid.UUID,
	memberEmails []string,
	optionIds []int32,
	createFn func(image *calendar.Image, members []calendar.Member, options []option.Option) (*calendar.Calendar, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// オプション取得
		options, err := psql.FindOptionsByIds(ctx, tx, optionIds)
		if err != nil {
			return err
		}
		// ユーザー取得
		users, err := psql.FindUsersByEmails(ctx, tx, memberEmails)
		if err != nil {
			return err
		}
		// メンバーリスト作成
		members := make([]calendar.Member, len(users))
		for i, user := range users {
			member, err := calendar.NewMember(user.Id, user.Name)
			if err != nil {
				continue
			}
			members[i] = *member
		}
		// 画像を取得
		image, err := psql.FindImageById(ctx, tx, imageId)
		if err != nil {
			return err
		}
		// スケジュール作成関数を実行
		calendar, err := createFn(image, members, options)
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
				"joinedAt": sql.NullTime{
					Time:  member.JoinedAt,
					Valid: !member.JoinedAt.IsZero(),
				},
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

func (r *CalendarPsqlRepository) Update(
	ctx context.Context,
	userId uuid.UUID,
	calendarId uuid.UUID,
	optionIds []int32,
	updateFn func(calendar *calendar.Calendar, options []option.Option) (*calendar.Calendar, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {

		// カレンダー取得
		existingCalendar, err := psql.FindCalendarByUserIdAndId(ctx, tx, userId, calendarId)
		if err != nil {
			return err
		}

		// オプション取得
		options, err := psql.FindOptionsByIds(ctx, tx, optionIds)
		if err != nil {
			return err
		}

		// 更新関数実行
		cal, err := updateFn(existingCalendar, options)
		if err != nil {
			return err
		}

		// カレンダー更新
		query := `
            UPDATE calendars SET
                name = :name,
                color = :color
            WHERE id = :id AND deleted_at IS NULL
        `
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":    cal.Id,
			"name":  cal.Name,
			"color": cal.Color,
		})
		if err != nil {
			return err
		}

		// 全置換
		// オプション全削除
		_, err = tx.ExecContext(ctx,
			"DELETE FROM calendar_options WHERE calendar_id = $1", cal.Id)
		if err != nil {
			return err
		}
		// オプション設定
		if len(cal.Options) > 0 {
			query = `
                INSERT INTO calendar_options(calendar_id, option_id)
                VALUES (:calendarId, :optionId)
            `
			calendarOptionMaps := []map[string]any{}
			for _, option := range cal.Options {
				calendarOptionMaps = append(calendarOptionMaps, map[string]any{
					"calendarId": cal.Id,
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
