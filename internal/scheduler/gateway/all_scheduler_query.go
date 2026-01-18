package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type AllSchedulerModel struct {
	Id         uuid.UUID `db:"id"`
	CalendarId uuid.UUID `db:"calendar_id"`
	UserId     uuid.UUID `db:"user_id"`
	Title      string    `db:"title"`
	Memo       string    `db:"memo"`
	LimitTime  time.Time `db:"limitTime"`
	IsAllDay   bool      `db:"is_allday"`
	IsDone     bool      `db:"is_done"`
}

func (g *SchedulerPsqlRepository) FindAllSchedulerById(ctx context.Context, calendarId, userId uuid.UUID) ([]scheduler.Scheduler, error) {
	sql := `
	SELECT s.id, s.calendar_id, s.user_id, s.title, s.memo, s.limit_time, s.is_allday, s.is_done,
	FROM scheduler s
	LEFT JOIN scheduler_members sm ON sm.user_id = $2
	WHERE s.calendarId = $1 AND sm.joined_at IS NULL
	`
	var rows []AllSchedulerModel
	err := g.db.SelectContext(ctx, &rows, sql, calendarId, userId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
