package query

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserQueryInput struct {
	UserId uuid.UUID
}

type UserQueryOutput struct {
	Id        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ユーザー情報を取得する
func (q *UserQuery) GetUser(ctx context.Context, input UserQueryInput) (*UserQueryOutput, error) {
	user, err := q.userRepository.FindById(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted() {
		return nil, errors.New("the user is deleted")
	}
	return &UserQueryOutput{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
