package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type UserStatusModel struct {
	UserId  uuid.UUID `db:"user_id"`
	Comment string    `db:"comment"`
	Date    time.Time `db:"date"`
	Status  int8      `db:"status"`
}

func (r *SchedulerPsqlRepository) FindStatusById(ctx context.Context, schedulerId, userId uuid.UUID) (*output.SchedulerUserOutput, error) {
	sql := `
	SELECT sa.user_id, sa.comment, ss.date, ss.status
	FROM scheduler_attendance sa
	LEFT JOIN scheduler s ON sa.scheduler_id = s.id
	LEFT JOIN scheduler_status ss ON sa.id = ss.attendance_id
	LEFT JOIN calendar_members cm ON s.calendar_id = cm.calendar_id
	WHERE sa.scheduler_id = $1 AND cm.user_id = $2 AND cm.joined_at IS NOT NULL
	`
	var rows []UserStatusModel
	err := r.db.SelectContext(ctx, &rows, sql, schedulerId, userId)
	if err != nil {
		return nil, err
	}
	statuses := make([]output.UserStatus, len(rows))
	for i, v := range rows {
		statuses[i] = output.UserStatus{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	result := output.SchedulerUserOutput{
		UserId: rows[0].UserId,
		Status: statuses,
	}
	return &result, nil
}
