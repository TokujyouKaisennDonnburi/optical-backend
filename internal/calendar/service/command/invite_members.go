package command

import (
	"context"
	"fmt"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type MemberCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Emails     []string
}

func (c *CalendarCommand) InviteMember(ctx context.Context, input MemberCreateInput) error {
	// emails validate
	emails := make([]user.Email, 0, len(input.Emails))
	for _, email := range input.Emails {
		validated, err := user.NewEmail(email)
		if err != nil {
			return err
		}
		emails = append(emails, validated)
	}
	err := c.memberRepository.Invite(ctx, input.UserId, input.CalendarId, emails)
	if err != nil {
		return err
	}
	go func() {
		content := getEmailContent(input.CalendarId)
		err := c.emailRepository.NotifyAll(
			ctx,
			"OptiCal: メンバーに招待されました",
			content,
			emails,
		)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

func getEmailContent(calendarId uuid.UUID) string {
	message := `
	カレンダーに招待されました。
	参加リンク：https://localhost:3000/calendars/%s/join
	`
	return fmt.Sprintf(message, calendarId.String())
}
