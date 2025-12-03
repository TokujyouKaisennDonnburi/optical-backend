package repository

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberRepository interface {
	Create(ctx context.Context, calendarId uuid.UUID, email string)(calendar.Member, error)
}
