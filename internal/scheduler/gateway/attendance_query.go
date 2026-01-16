package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type SchedulerModel struct {
	Id         uuid.UUID `db:"id"`
	CalendarId uuid.UUID `db:"calendar_id"`
	UserId     uuid.UUID `db:"user_id"`
	Title      string    `db:"title"`
	Memo       string    `db:"memo"`
	LimitTime  time.Time `db:"limit_time"`
	IsAllDay   bool      `db:"is_allday"`
	Date       time.Time `db:"date"`
	StartTime  time.Time `db:"start_time"`
	EndTime    time.Time `db:"end_time"`
}

func (r *SchedulerPsqlRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*output.SchedulerAttendanceQuery, error) {
	query := `
	SELECT s.id, s.calendar_id, s.user_id, s.title, s.memo, s.limit_time, s.is_allday,
	pd.date, pd.start_time, pd.end_time
	FROM scheduler s
	LEFT JOIN scheduler_possible_date pd ON pd.scheduler_id = s.id
	WHERE s.id = $1
	`
	var row []SchedulerModel
	err := r.db.SelectContext(ctx, &row, query, id)
	if err != nil {
		return nil, err
	}
	if len(row) == 0 {
		return nil, errors.New("scheduler not found")
	}
	dates := make([]output.PossibleDateOutput, len(row))
	for i, v := range row {
		dates[i] = output.PossibleDateOutput{
			Date:      v.Date,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		}
	}
	result := output.SchedulerAttendanceQuery{
		Id:           row[0].Id,
		CalendarId:   row[0].CalendarId,
		UserId:       row[0].UserId,
		Title:        row[0].Title,
		Memo:         row[0].Memo,
		LimitTime:    row[0].LimitTime,
		IsAllDay:     row[0].IsAllDay,
		PossibleDate: dates,
	}
	return &result, nil
}
