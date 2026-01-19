package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"gopkg.in/mail.v2"
)

type GmailRepository struct {
	dialer *mail.Dialer
	email  string
}

func NewGmailRepository(dialer *mail.Dialer, email string) *GmailRepository {
	if dialer == nil {
		panic("email dialer is nil")
	}
	if email == "" {
		panic("email address is empty")
	}
	return &GmailRepository{
		dialer: dialer,
		email:  email,
	}
}

func (r *GmailRepository) NotifyAll(ctx context.Context, subject, content string, emails []user.Email) error {
	m := mail.NewMessage()
	m.SetHeader("From", r.email)
	m.SetHeader("To", r.email)
	mails := make([]string, len(emails))
	for i, email := range emails {
		mails[i] = string(email)
	}
	m.SetHeader("Bcc", mails...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)
	return r.dialer.DialAndSend(m)
}
