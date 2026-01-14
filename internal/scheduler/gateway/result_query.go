package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type ResultModel struct {
	OwnerId   uuid.UUID `db:"owner_id"`
	Title     string    `db:"title"`
	Memo      string    `db:"memo"`
	LimitTime time.Time `db:"limit_time"`
	IsAllDay  bool      `db:"is_allDay"`
	UserId    uuid.UUID `db:"user_id"`
	UserName  uuid.UUID `db:"user_name"`
}
type DateModel struct {
	Date      time.Time `db:"date"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

func (r *SchedulerPsqlRepository) ResultGateway(
	ctx context.Context,
	calendarId, schedulerId, userId uuid.UUID,
) (*output.SchedulerResultOutput, error) {
	// result
	sql := `
	SELECT s.user_id AS owner_id, s.title, s.memo, s.limit_time, s.is_allday,
	cm.user_id, cm.joined_at,
	u.name AS user_name
	FROM scheduler s
	LEFT JOIN calendar_members cm ON cm.calendar_id = $1
	LEFT JOIN users u ON u.id = cm.id
	WHERE s.id = $2
		`
	var resultModels []ResultModel
	err := r.db.SelectContext(ctx, &resultModels, sql, calendarId, schedulerId)
	if err != nil {
		return nil, err
	}
	if len(resultModels) == 0 {
		return nil, errors.New("scheduler not found")
	}
	members := make([]output.MemberOutput, len(resultModels))
	for i, v := range resultModels {
		members[i] = output.MemberOutput{
			UserId:   v.UserId,
			UserName: v.UserName,
		}
	}

	// date
	sql = `
	SELECT pd.date, pd.start_time, pd.end_time
	FROM scheduler_possible_date pd ON pd.scheduler_id = $1
	`
	var dateModels []DateModel
	err = r.db.SelectContext(ctx, &dateModels, sql, schedulerId)
	if err != nil {
		return nil, err
	}
	dates := make([]output.DateOutput, len(dateModels))
	for i, v := range dateModels {
		dates[i] = output.DateOutput{
			Date:      v.Date,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		}
	}
	output := output.SchedulerResultOutput{
		OwnerId:   resultModels[0].OwnerId,
		Title:     resultModels[0].Title,
		Memo:      resultModels[0].Memo,
		LimitTime: resultModels[0].LimitTime,
		IsAllDay:  resultModels[0].IsAllDay,
		Members:   members,
		Date:      dates,
	}
	return &output, nil
}
