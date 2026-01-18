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
	LimitTime  time.Time `db:"limit_time"`
	IsAllDay   bool      `db:"is_allday"`
	IsDone     bool      `db:"is_done"`
}

func (d *SchedulerPsqlRepository) FindAllSchedulerById(ctx context.Context, calendarId uuid.UUID) ([]scheduler.Scheduler, error) {
	sql := `
	SELECT s.id, s.calendar_id, s.user_id, s.title, s.memo, s.limit_time, s.is_allday, s.is_done
	FROM scheduler s
	WHERE s.calendar_id = $1
	`
	var rows []AllSchedulerModel
	err := d.db.SelectContext(ctx, &rows, sql, calendarId)
	if err != nil {
		return nil, err
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
	return result, nil
}
