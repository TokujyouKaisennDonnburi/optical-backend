package command

import (
	"context"
	"fmt"
	"os"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MemberCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Emails     []string
}

// 招待を作成
func (c *CalendarCommand) CreateInvitations(ctx context.Context, input MemberCreateInput) error {
	// 2. Validate emails
	validEmails := make([]string, 0, len(input.Emails))
	for _, email := range input.Emails {
		_, err := user.NewEmail(email)
		if err != nil {
			return err
		}
		validEmails = append(validEmails, email)
	}

	// 3. Create Invitations
	invitations := make([]*calendar.Invitation, 0, len(validEmails))
	for _, email := range validEmails {
		inv, err := calendar.NewInvitation(input.CalendarId, email)
		if err != nil {
			return err
		}
		invitations = append(invitations, inv)
	}

	// 4. Save to DB
	err := c.memberRepository.CreateInvitations(ctx, invitations)
	if err != nil {
		return err
	}

	// 5. Send Emails
	go func() {
		baseUrl := getFontendUrl()
		for _, inv := range invitations {
			content := getInvitationEmailContent(baseUrl, input.CalendarId, inv.Token)
			// Converting string to user.Email for the repository interface
			userEmail, _ := user.NewEmail(inv.Email)
			err := c.emailRepository.NotifyAll(
				context.Background(), // Use a new context for background task
				"OptiCal: メンバーに招待されました",
				content,
				[]user.Email{userEmail},
			)
			if err != nil {
				logrus.WithError(err).WithField("email", inv.Email).Error("failed to send invitation email")
			}
		}
	}()

	return nil
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
		content := getEmailContent(getFontendUrl(), input.CalendarId)
		err := c.emailRepository.NotifyAll(
			ctx,
			"OptiCal: メンバーに招待されました",
			content,
			emails,
		)
		if err != nil {
			logrus.WithError(err).Error("failed to send emails")
		}
	}()
	return nil
}

func getFontendUrl() string {
	baseUrl := os.Getenv("FRONTEND_BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:3000"
	}
	return baseUrl
}

func getEmailContent(baseUrl string, calendarId uuid.UUID) string {
	message := `
	カレンダーに招待されました。
	参加リンク：%s/calendars/%s/join
	`
	return fmt.Sprintf(message, baseUrl, calendarId.String())
}

func getInvitationEmailContent(baseUrl string, calendarId, token uuid.UUID) string {
	message := `
	カレンダーに招待されました。
	以下のリンクから参加してください（リンクは30日間有効です）。
	
	参加リンク：%s/calendars/%s/join?token=%s
	`
	return fmt.Sprintf(message, baseUrl, calendarId.String(), token.String())
}
