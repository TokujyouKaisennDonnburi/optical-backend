package command

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberCreateInput struct {
	UserId 		uuid.UUID
	CalendarId  uuid.UUID
	Email 		string 
}

func (c *CalendarCommand) CreateMember (ctx context.Context, input MemberCreateInput)error {
	err := c.memberRepository.Create(ctx, input.UserId, input.CalendarId, input.Email)
	if err != nil {
		return err
	}
	return nil
}
