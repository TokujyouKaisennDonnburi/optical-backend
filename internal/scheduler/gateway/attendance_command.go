package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *SchedulerPsqlRepository) AddAttendance(
	ctx context.Context,
	id, schedulerId, userId uuid.UUID,
	comment string,
	schedulerStatus []scheduler.SchedulerStatus,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// attendance
		query := `
			INSERT INTO scheduler_attendance(id, scheduler_id, user_id, comment)
			VALUES (:id, :schedulerId, :userId, :comment)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":          id,
			"schedulerId": schedulerId,
			"userId":      userId,
			"comment":     comment,
		})
		if err != nil {
			return err
		}

		// status
		query = `
			INSERT INTO scheduler_status(attendance_id, date, status)
			VALUES (:attendanceId, :date, :status)
		`
		for _, v := range schedulerStatus {
			_, err = tx.NamedExecContext(ctx, query, map[string]any{
				"attendanceId": id,
				"date":         v.Date,
				"status":       v.Status,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
