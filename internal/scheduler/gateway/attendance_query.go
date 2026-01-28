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

func (r *SchedulerPsqlRepository) FindAttendanceById(ctx context.Context, calendarId, schedulerId, userId uuid.UUID) ([]output.SchedulerAttendanceOutput, error) {
	sql := `
	SELECT sa.user_id, sa.comment, ss.date, ss.status
	FROM scheduler_attendance sa
	LEFT JOIN scheduler_status ss ON ss.attendance_id = sa.id
	INNER JOIN scheduler s ON s.id = sa.scheduler_id AND s.calendar_id = $1
	INNER JOIN calendar_members cm_att ON cm_att.calendar_id = $1 AND cm_att.user_id = sa.user_id AND cm_att.joined_at IS NOT NULL
	WHERE sa.scheduler_id = $2
	  AND EXISTS (
		SELECT 1 FROM calendar_members cm_req
		WHERE cm_req.calendar_id = $1 AND cm_req.user_id = $3 AND cm_req.joined_at IS NOT NULL
	  )
	ORDER BY sa.user_id, ss.date
	`
	var rows []UserStatusModel
	err := r.db.SelectContext(ctx, &rows, sql, calendarId, schedulerId, userId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []output.SchedulerAttendanceOutput{}, nil
	}
	results := make([]output.SchedulerAttendanceOutput, 0)
	indexByUser := make(map[uuid.UUID]int)
	for _, v := range rows {
		idx, ok := indexByUser[v.UserId]
		if !ok {
			results = append(results, output.SchedulerAttendanceOutput{
				UserId:  v.UserId,
				Comment: v.Comment,
				Status:  []output.StatusOutput{},
			})
			idx = len(results) - 1
			indexByUser[v.UserId] = idx
		}
		if !v.Date.IsZero() {
			results[idx].Status = append(results[idx].Status, output.StatusOutput{
				Date:   v.Date,
				Status: v.Status,
			})
		}
	}
	return results, nil
}
