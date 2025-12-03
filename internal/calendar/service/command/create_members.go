package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberCreateInput struct {
	UserId 		uuid.UUID
	CalendarId  uuid.UUID
	Email 		string 
}

func (c *CalendarCommand) CreateMember (ctx context.Context, input MemberCreateInput) (*calendar.Member, error) {
	var newMember *calendar.Member
	_, err := c.memberRepository.Create(input.Email)
	if err != nil {
		return nil, err 
	}
	return newMember, nil
}
