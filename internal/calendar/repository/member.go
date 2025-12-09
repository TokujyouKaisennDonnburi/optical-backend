package repository

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberRepository interface {
	Create(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email)error
}
