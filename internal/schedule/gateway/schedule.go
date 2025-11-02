package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	optionGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/option/gateway"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SchedulePsqlRepository struct {
	db *sqlx.DB
}

func NewSchedulePsqlRepository(db *sqlx.DB) *SchedulePsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &SchedulePsqlRepository{
		db: db,
	}
}

// スケジュールを新規作成する
func (r *SchedulePsqlRepository) Create(
	ctx context.Context,
	userId uuid.UUID,
	optionIds []uuid.UUID,
	createFn func(user *user.User, options []option.Option) (*schedule.Schedule, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// オプション取得
		options, err := optionGateway.FindOptionsByIds(ctx, tx, optionIds)
		if err != nil {
			return err
		}
		// ユーザー取得
		user, err := userGateway.FindUserById(ctx, tx, userId)
		if err != nil {
			return err
		}
		// スケジュール作成関数を実行
		schedule, err := createFn(user, options)
		if err != nil {
			return err
		}
		// スケジュール作成
		query := `
			INSERT INTO schedules(id, name)
			VALUES (:id, :name)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":   schedule.Id,
			"name": schedule.Name,
		})
		if err != nil {
			return err
		}
		// メンバー作成
		query = `
			INSERT INTO schedule_members(schedule_id, user_id, joined_at)
			VALUES (:scheduleId, :userId, :joinedAt)
		`
		scheduleMemberMaps := []map[string]any{}
		for _, member := range schedule.Members {
			scheduleMemberMaps = append(scheduleMemberMaps, map[string]any{
				"scheduleId": schedule.Id,
				"userId":     member.UserId,
				"joinedAt":   member.JoinedAt,
			})
		}
		_, err = tx.NamedExecContext(ctx, query, scheduleMemberMaps)
		if err != nil {
			return err
		}
		// オプション設定
		if len(options) > 0 {
			query = `
				INSERT INTO schedule_options(schedule_id, option_id)
				VALUES (:scheduleId, :optionId)
			`
			scheduleOptionMaps := []map[string]any{}
			for _, option := range schedule.Options {
				scheduleMemberMaps = append(scheduleMemberMaps, map[string]any{
					"scheduleId": schedule.Id,
					"optionId":   option.Id,
				})
			}
			_, err = tx.NamedExecContext(ctx, query, scheduleOptionMaps)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
