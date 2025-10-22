package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
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
func (r *SchedulePsqlRepository) Create(ctx context.Context, schedule *schedule.Schedule) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// スケジュール作成
		query := `
			INSERT INTO schedules(id, name)
			VALUES (:id, :name)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
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
		// オプション作成
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
		return nil
	})
}
