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

func (r *SchedulerPsqlRepository) FindStatusById(ctx context.Context, schedulerId, userId uuid.UUID) (output.SchedulerUserOutput, error) {
	sql := `
	SELECT sa.user_id, sa.comment
	FROM scheduler_attendance sa
	LEFT JOIN scheduler_status ss ON sa.id = ss.attendance_id
	WHERE
	`
}
