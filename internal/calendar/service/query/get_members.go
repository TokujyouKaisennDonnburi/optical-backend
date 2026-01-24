package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

// handlerから受け取るデータ
type MemberQueryInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

// メンバー一覧取得
func (q *MemberQuery) GetMembers(ctx context.Context, input MemberQueryInput) ([]output.MembersQueryOutput, error) {
	// 権限チェック: ユーザーがカレンダーのメンバーか確認
	exists, err := q.memberRepository.ExistsMemberByUserIdAndCalendarId(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperr.ForbiddenError("not a member of this calendar")
	}

	// メンバー一覧取得
	return q.memberRepository.FindMembers(ctx, input.CalendarId)
}
