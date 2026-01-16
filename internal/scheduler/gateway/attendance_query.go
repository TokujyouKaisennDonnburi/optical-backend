package gateway

import (
	"context"
	"errors"
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

func (r *SchedulerPsqlRepository) FindAttendanceById(ctx context.Context, calendarId, schedulerId, userId uuid.UUID) (*output.SchedulerAttendanceOutput, error) {
	sql := `
	SELECT sa.user_id, sa.comment, ss.date, ss.status
	FROM scheduler_attendance sa
	LEFT JOIN scheduler_status ss ON ss.attendance_id = sa.id
	INNER JOIN calendar_members cm ON cm.calendar_id = $1
	WHERE sa.scheduler_id = $2 AND cm.user_id = $3 AND cm.joined_at IS NOT NULL
	`
	var rows []UserStatusModel
	err := r.db.SelectContext(ctx, &rows, sql, calendarId, schedulerId, userId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("scheduler status is not found")
	}
	statuses := make([]output.StatusOutput, len(rows))
	for i, v := range rows {
		statuses[i] = output.StatusOutput{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	result := output.SchedulerAttendanceOutput{
		UserId:  rows[0].UserId,
		Comment: rows[0].Comment,
		Status:  statuses,
	}
	return &result, nil
}
