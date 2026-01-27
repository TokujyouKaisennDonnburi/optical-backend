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
	// メンバー一覧取得
	members, err := q.memberRepository.FindMembers(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}

	// 権限チェック: ユーザーがカレンダーのメンバーか確認
	isMember := false
	for _, m := range members {
		if m.UserId == input.UserId {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, apperr.ForbiddenError("not a member of this calendar")
	}

	return members, nil
}
