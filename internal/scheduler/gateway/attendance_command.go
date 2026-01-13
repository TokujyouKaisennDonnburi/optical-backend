package gateway

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

func (r *SchedulerPsqlRepository) AddAttendance(
	id, calendarId, userId uuid.UUID,
	comment string,
	schedulerStatus []scheduler.SchedulerStatus,
) error {
	query := `
		INSERT INTO scheduler_attendance(id, scheduler_id, user_id, comment)
		VALUES (:id, :schedulerId, :userId, comment)
	`
	query = `
		INSERT INTO scheduler_status(attendance_id, date, status)
		VALUES (:id, :date, :status)
	`
}
