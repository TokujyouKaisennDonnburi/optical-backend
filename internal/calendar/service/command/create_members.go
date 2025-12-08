package command

import (
	"context"

	"github.com/google/uuid"
)

type MemberCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Emails     []string
}

func (c *CalendarCommand) CreateMember(ctx context.Context, input MemberCreateInput) error {
	err := c.memberRepository.Create(ctx, input.UserId, input.CalendarId, input.Emails[])
	if err != nil {
		return err
	}
	return nil
}
