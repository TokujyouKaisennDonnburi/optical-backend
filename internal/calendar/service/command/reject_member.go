package command

import (
	"context"

	"github.com/google/uuid"
)

type RejectMemberInput struct {
	UserId uuid.UUID
	CalendarId uuid.UUID
}

func (c *CalendarCommand)RejectMember(ctx context.Context, input RejectMemberInput)error{
	return c.memberRepository.Reject(ctx, input.UserId, input.CalendarId)
}
