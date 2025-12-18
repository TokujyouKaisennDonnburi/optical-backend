package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"gopkg.in/mail.v2"
)

type GmailRepository struct {
	dialer *mail.Dialer
}

func NewGmailRepository(dialer *mail.Dialer) *GmailRepository {
	if dialer == nil {
		panic("email dialer is nil")
	}
	return &GmailRepository{
		dialer: dialer,
	}
}

func (r *GmailRepository) NotifyAll(ctx context.Context, subject, content string, emails []user.Email) error {
	m := mail.NewMessage()
	m.SetHeader("From", r.dialer.Username)
	m.SetHeader("To", r.dialer.Username)
	mails := make([]string, len(emails))
	for i, email := range emails {
		mails[i] = string(email)
	}
	m.SetHeader("Bcc", mails...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)
	return r.dialer.DialAndSend(m)
}
