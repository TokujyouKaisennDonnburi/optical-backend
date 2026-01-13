package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type AttendanceInput struct {
	CalendarId  uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
	Comment     string
	Status      []StatusInput
}

type StatusInput struct {
	Date   time.Time
	Status int8
}

func (c *SchedulerCommand) AddAttendanceCommand(ctx context.Context, input AttendanceInput) error {
	// option check
	options, err := c.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return err
	}
	hasOption := false
	for _, x := range options {
		if x.Id == option.OPTION_SCHEDULER {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return apperr.ForbiddenError("option not enabled")
	}
	// domain
	domain, err := scheduler.NewAttendance(input.SchedulerId, input.UserId, input.Comment)
	if err != nil {
		return err
	}
	domains := make([]scheduler.SchedulerStatus, 0, len(input.Status))
	for _, v := range input.Status {
		status, err := scheduler.NewStatus(v.Date, v.Status)
		if err != nil {
			return err
		}
		domains = append(domains, *status)
	}
	err = c.schedulerRepository.AddAttendance(ctx, domain.Id, domain.SchedulerId, domain.UserId, domain.Comment,domains)
	if err != nil {
		return err
	}
	return nil
}
