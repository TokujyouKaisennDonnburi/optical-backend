package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query/output"
	"github.com/google/uuid"
)

type ResultModel struct {
	UserId    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	Memo      string    `db:"memo"`
	LimitTime time.Time `db:"limit_time"`
	IsAllDay  bool      `db:"is_allDay"`
	Date      time.Time `db:"date"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}
type MemberModel struct {
	UserId   uuid.UUID
	UserName uuid.UUID
}
func (d *SchedulerPsqlRepository)ResultGateway(ctx context.Context, calendarId, schedulerId,userId uuid.UUID)(output.SchedulerResultOutput, error){
	sql := SELECT scheduler.
}
