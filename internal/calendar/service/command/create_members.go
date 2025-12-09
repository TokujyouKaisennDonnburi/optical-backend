package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type MemberCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Emails     []string
}

func (c *CalendarCommand) CreateMember(ctx context.Context, input MemberCreateInput) error {
	// emails validate
	emails := make([]string, 0, len(input.Emails))
	for _, email := range input.Emails {
		validated, err := user.NewEmail(email)
		if err != nil {
			return err
		}
		emails = append(emails, string(validated))
	}
	err := c.memberRepository.Create(ctx, input.UserId, input.CalendarId, input.Emails)
	if err != nil {
		return err
	}
	return nil
}
