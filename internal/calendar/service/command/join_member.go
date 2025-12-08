package command

import (
	"context"

	"github.com/google/uuid"
)

type CalendarJoinInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

func (c *CalendarCommand) JoinMember(ctx context.Context, input CalendarJoinInput)error{
	err := c.memberRepository.Join(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return err
	}
	return nil
}
