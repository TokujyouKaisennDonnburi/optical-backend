package command

import (
	"context"

	"github.com/google/uuid"
)

type EventDeleteInput struct {
	EventId uuid.UUID
	UserId  uuid.UUID
}

func (c *EventCommand) Delete(ctx context.Context, input EventDeleteInput) error {
	err := c.eventRepository.Delete(ctx, input.EventId, input.UserId)
	if err != nil {
		return err
	}
	return nil
}
