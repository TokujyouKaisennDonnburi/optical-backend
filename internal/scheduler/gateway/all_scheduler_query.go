package gateway

import (
	"context"
	"errors"
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
	LimitTime  time.Time `db:"limit_time"`
	IsAllDay   bool      `db:"is_allday"`
	IsDone     bool      `db:"is_done"`
}

func (d *SchedulerPsqlRepository) FindAllSchedulerById(ctx context.Context, calendarId, userId uuid.UUID) (*scheduler.Scheduler, error) {
	sql := `
	SELECT s.id, s.calendar_id, s.user_id, s.title, s.memo, s.limit_time, s.is_allday, s.is_done,
	FROM scheduler s
	LEFT JOIN calendar_members cm ON cm.user_id = $2
	WHERE s.calendar_id = $1
	`
	var rows []AllSchedulerModel
	err := d.db.SelectContext(ctx, &rows, sql, calendarId, userId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("scheduler not found")
	}
	result := make([]scheduler.Scheduler, len(rows))
	for i, v := range rows {
		result[i] = scheduler.Scheduler{
			Id:         v.Id,
			CalendarId: v.CalendarId,
			UserId:     v.UserId,
			Title:      v.Title,
			Memo:       v.Memo,
			LimitTime:  v.LimitTime,
			IsAllDay:   v.IsAllDay,
			IsDone:     v.IsDone,
		}
	}
	return *result, nil
}
