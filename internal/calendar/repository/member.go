package repository

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberRepository interface {
	Create(ctx context.Context, userId, calendarId uuid.UUID, emails []string)error
}
