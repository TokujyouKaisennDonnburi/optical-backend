package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type CalendarJoinInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type JoinWithTokenInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Token      uuid.UUID
}

func (c *CalendarCommand) JoinMember(ctx context.Context, input CalendarJoinInput) error {
	err := c.memberRepository.Join(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return err
	}
	return nil
}

// 個別トークンからカレンダーメンバーに追加
func (c *CalendarCommand) JoinMemberWithToken(ctx context.Context, input JoinWithTokenInput) error {
	return c.transactor.Transact(ctx, func(ctx context.Context) error {
		// トークン検証
		invitation, err := c.memberRepository.FindInvitationByToken(ctx, input.Token)
		if err != nil {
			return err
		}
		if invitation.CalendarId != input.CalendarId {
			return apperr.ValidationError("invalid token for this calendar")
		}
		if invitation.IsExpired() {
			return apperr.ValidationError("invitation expired")
		}
		if invitation.IsUsed() {
			return apperr.ValidationError("invitation already used")
		}

		// 招待使用済み更新
		invitation.MarkAsUsed(input.UserId)
		if err := c.memberRepository.MarkInvitationAsUsed(ctx, invitation); err != nil {
			return err
		}

		// メンバー追加
		return c.memberRepository.AddMember(ctx, input.CalendarId, input.UserId)
	})
}
