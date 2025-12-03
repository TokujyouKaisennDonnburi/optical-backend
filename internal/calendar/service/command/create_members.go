package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"golang.org/x/net/context"
)

type MemberCreateInput struct {
	Email string 
}

func (c *CalendarCommand) CreateMember (ctx context.Context, in CalendarMemberInput) (calendar.Member) {
	var newMember *calendar.Member
	calendar , err := c.memberRepository.Create(MemberCreateInput.Email)
	if err != nil {
		return nil, err 
	}
	return newMember, nil
}
