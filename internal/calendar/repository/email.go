package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

type EmailRepository interface {
	NotifyAll(ctx context.Context, subject, content string, emails []user.Email) error
}
