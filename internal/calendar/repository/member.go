package repository

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type MemberRepository interface {
	// カレンダーメンバー取得
	FindMembers(ctx context.Context, calendarId uuid.UUID) ([]output.MembersQueryOutput, error)
	// 既存のメソッド(calendar_members操作)
	Invite(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error
	Join(ctx context.Context, userId, calendarId uuid.UUID) error
	Reject(ctx context.Context, userId, calendarId uuid.UUID) error

	// calendar_memberに登録
	AddMember(ctx context.Context, calendarId, userId uuid.UUID) error

	// 招待関連(calendar_invitations操作)
	// 招待を一括作成
	CreateInvitations(ctx context.Context, invitations []*calendar.Invitation) error
	// トークンで招待を取得
	FindInvitationByToken(ctx context.Context, token uuid.UUID) (*calendar.Invitation, error)
	// カレンダーIDで招待一覧を取得(未使用かつ有効期限内)
	FindPendingInvitationsByCalendarId(ctx context.Context, calendarId uuid.UUID) ([]*calendar.Invitation, error)
	// 招待を使用済みにする
	MarkInvitationAsUsed(ctx context.Context, invitation *calendar.Invitation) error
}
