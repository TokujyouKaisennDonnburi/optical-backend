package command

import (
	"context"

	"github.com/google/uuid"
)

type CalendarDeleteInput struct {
	CalendarId uuid.UUID
	UserId     uuid.UUID
}

func (c *CalendarCommand) DeleteCalendar(ctx context.Context, input CalendarDeleteInput) error {
	err := c.calendarRepository.Delete(ctx, input.CalendarId, input.UserId)
	if err != nil {
		return err
	}
	return nil
}
