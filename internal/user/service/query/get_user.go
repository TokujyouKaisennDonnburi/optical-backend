package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query/output"
	"github.com/google/uuid"
)

type UserQueryInput struct {
	UserId uuid.UUID
}

// ユーザー情報を取得する
func (q *UserQuery) GetUser(
	ctx context.Context,
	input UserQueryInput,
) (*output.UserQueryOutput, error) {
	return q.userRepository.FindProfileById(ctx, input.UserId)
}
