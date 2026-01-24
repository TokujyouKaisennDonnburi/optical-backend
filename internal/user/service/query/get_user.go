package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
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
	output, err := q.userRepository.FindProfileById(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	if output.Avatar.Valid {
		if output.Avatar.IsRelativePath {
			output.AvatarUrl = storage.GetImageStorageBaseUrl() + "/" + output.Avatar.Url
		} else {
			output.AvatarUrl = output.Avatar.Url
		}
	}
	return output, nil
}
