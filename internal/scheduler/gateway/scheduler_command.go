package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SchedulerPsqlRepository struct {
	db *sqlx.DB
}

func NewSchedulerPsqlRepository(db *sqlx.DB) *SchedulerPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &SchedulerPsqlRepository{
		db: db,
	}
}

// create
func (r *SchedulerPsqlRepository) CreateScheduler(
	ctx context.Context,
	id, calendarId, userId uuid.UUID,
	title, memo string,
	possibleDates []scheduler.PossibleDate,
	limitTime time.Time,
	isAllDay bool,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// scheduler
		sql := `
		INSERT INTO scheduler(id, calendar_id, user_id, title, memo, limit_time, is_allday)
		VALUES(:id, :calendarId, :userId, :title, :memo, :limitTime, :isAllDay)
	`
		_, err := tx.NamedExecContext(ctx, sql, map[string]any{
			"id":         id,
			"calendarId": calendarId,
			"userId":     userId,
			"title":      title,
			"memo":       memo,
			"limitTime":  limitTime,
			"isAllDay":   isAllDay,
		})
		if err != nil {
			return err
		}

		// possibleDate
		sql = `
		INSERT INTO scheduler_possible_date(scheduler_id, date, start_time, end_time)
		VALUES(:schedulerId, :date, :startTime, :endTime)
	`
		for _, db := range possibleDates {
			_, err = tx.NamedExecContext(ctx, sql, map[string]any{
				"schedulerId": id,
				"date":        db.Date,
				"startTime":   db.StartTime,
				"endTime":     db.EndTime,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
