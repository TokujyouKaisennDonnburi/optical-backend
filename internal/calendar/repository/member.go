package repository

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberRepository interface {
	Invite(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error
	Join(ctx context.Context, userId, calendarId uuid.UUID) error
	Reject(ctx context.Context, userId, calendarId uuid.UUID) error

	// 権限チェック
	ExistsMemberByUserIdAndCalendarId(ctx context.Context, userId, calendarId uuid.UUID) (bool, error)
	// メンバー
	FindMembers(ctx context.Context, calendarId uuid.UUID) ([]output.MembersQueryOutput, error)
}
