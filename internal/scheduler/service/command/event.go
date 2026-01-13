package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type SchedulerEventInput struct {
	CalendarId  uuid.UUID
	SchedulerId uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
}

func (c *SchedulerCommand) SchedulerEvent(ctx context.Context, input SchedulerEventInput)error{
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
}
